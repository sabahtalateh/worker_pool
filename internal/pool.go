package internal

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type workerList struct {
	mu      sync.Mutex
	total   int64
	working int64
}

func (wl *workerList) check(
	ctx context.Context,
	scheduler *Scheduler,
	wg *sync.WaitGroup,
	db *sqlx.DB,
) {
	println("check")

	wl.mu.Lock()
	defer wl.mu.Unlock()

	fmt.Printf("%d workers to run\n", wl.total-wl.working)

	for i := wl.working; i < wl.total; i++ {
		wg.Add(1)
		wl.working++
		go wl.runWorker(ctx, scheduler, wg, db)
	}
}

func (wl *workerList) runWorker(
	ctx context.Context,
	scheduler *Scheduler,
	wg *sync.WaitGroup,
	db *sqlx.DB,
) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover from: %s\n", r)
			wl.mu.Lock()
			defer wl.mu.Unlock()
			wg.Done()
			wl.working--
		}
	}()

	w := NewWorker(scheduler, db)
	w.work(ctx, wg)
}

type Pool struct {
	checkEvery time.Duration
	workers    *workerList
	db         *sqlx.DB
	started    bool
	stop       bool
	shutdownWG sync.WaitGroup
	cancel     context.CancelFunc
	scheduler  *Scheduler
}

func NewPool(size int64, checkEvery time.Duration, db *sqlx.DB) (*Pool, error) {
	if size < 0 {
		return nil, errors.New("pool size must be gte 0")
	}

	return &Pool{
		checkEvery: checkEvery,
		workers:    &workerList{total: size, working: 0},
		stop:       false,
		db:         db,
		scheduler:  NewScheduler(db),
	}, nil
}

func (p *Pool) Run(ctx context.Context) error {
	if p.started {
		return errors.New("pool already started")
	}
	p.started = true
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	p.workers.check(ctx, p.scheduler, &p.shutdownWG, p.db)

	t := time.NewTicker(p.checkEvery)

	for {
		select {
		case <-t.C:
			p.workers.check(ctx, p.scheduler, &p.shutdownWG, p.db)
		case <-ctx.Done():
			t.Stop()
			return nil
		}
	}
}

func (p *Pool) Stop() {
	fmt.Printf("wait %d workers stop\n", p.workers.working)
	p.stop = true
	p.cancel()
	p.shutdownWG.Wait()
}
