package internal

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Scheduler struct {
	db *sqlx.DB
}

func NewScheduler(db *sqlx.DB) *Scheduler {
	return &Scheduler{db: db}
}

// Шедулер локает таблицу и выбирает подходящий id, какая то кастомная логика балансировки задач может быть тут
func (b *Scheduler) getTask(ctx context.Context, tx *sqlx.Tx) (bool, int64, error) {
	var id int64
	if err := tx.QueryRowContext(ctx, selectTaskId).Scan(&id); err != nil {
		return false, 0, err
	}

	// После выбора задачи шедулер помечает её как зашедуленую и больше её не трогает

	_ = tx.Commit()

	return true, id, nil
}
