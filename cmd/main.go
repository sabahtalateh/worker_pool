package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //nolint
	"github.com/sabahtalateh/worker_pool/internal"
)

func main() {
	ctx := context.Background()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	db, _ := sqlx.Connect("postgres", "host=127.0.0.1 user=postgres password=postgres dbname=worker_pool port=23345 sslmode=disable binary_parameters=yes")

	js := internal.NewJobService(db)
	_ = js.SubmitJob(ctx, "standard deployment", 1, 3)

	pool, _ := internal.NewPool(2, 5*time.Second, db)
	go func() {
		_ = pool.Run(ctx)
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Обработка Ctrl+C
	go func() {
		for sig := range signals {
			// sig is a ^C, handle it
			println(sig)
			pool.Stop()
			wg.Done()
		}
	}()

	go func() {
		time.Sleep(60 * time.Second)
		pool.Stop()
		wg.Done()
	}()

	wg.Wait()
}
