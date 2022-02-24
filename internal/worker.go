package internal

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Worker struct {
	id        string
	db        *sqlx.DB
	scheduler *Scheduler
}

func NewWorker(scheduler *Scheduler, db *sqlx.DB) *Worker {
	return &Worker{id: uuid.New().String(), scheduler: scheduler, db: db}
}

func (w *Worker) work(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			println("done " + w.id)
			return
		default:
			r := rand.Intn(100)
			if r > 50 {
				panic(fmt.Sprintf("oh it's %d", r))
			}

			tx, _ := w.db.Beginx()

			// Смотрим на тип таски и вызываем нужный процессор
			ok, taskId, _ := w.scheduler.getTask(ctx, tx)

			// Процессор возвращает результат, воркер обновляет статусы тасок

			if !ok {
				continue
			}

			fmt.Println(taskId)

			_ = tx.Commit()

			fmt.Printf("worker %s works on %d\n", w.id, r)
			time.Sleep(1 * time.Second)
		}
	}

}
