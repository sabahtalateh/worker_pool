package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type jobTemplate struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type templateNodes struct {
	ID          int64   `db:"id"`
	Root        bool    `db:"root"`
	Name        string  `db:"name"`
	TemplateID  int64   `db:"template_id"`
	Type        string  `db:"type"`
	ApplyTo     *string `db:"apply_to"`
	Composition *string `db:"composition"`
	G1          *int64  `db:"g1"`
	G2          *int64  `db:"g2"`
	TaskTypeID  *int64  `db:"task_type_id"`
}

type newTask struct {
	TypeID    int64  `db:"type_id"`
	Status    string `db:"status"`
	ServiceID *int64 `db:"service_id"`
}

type JobService struct {
	db *sqlx.DB
}

func NewJobService(db *sqlx.DB) *JobService {
	return &JobService{db: db}
}

func (s *JobService) SubmitJob(ctx context.Context, jobTemplateName string, servicesIds ...int64) error {
	// на основании шаблона и сервисов генерируются таски
	rr, err := s.db.QueryxContext(ctx, selectJobTemplate, jobTemplateName)
	defer rr.Close()

	if err != nil {
		return err
	}

	var template jobTemplate
	if rr.Rows.Next() {
		if err = rr.StructScan(&template); err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("no job `%s` exists", jobTemplateName))
	}

	rr, err = s.db.QueryxContext(ctx, selectJobTemplateNodes, template.ID)
	defer rr.Close()
	nodes := make([]templateNodes, 0)
	for rr.Rows.Next() {
		var node templateNodes
		if err = rr.StructScan(&node); err != nil {
			return err
		}
		nodes = append(nodes, node)
	}

	for _, node := range nodes {
		if *node.ApplyTo == "every_service" {
			for _, ID := range servicesIds {
				query, params, err := s.db.BindNamed(
					insertTask,
					newTask{TypeID: *node.TaskTypeID, Status: "new", ServiceID: &ID},
				)

				if err != nil {
					return err
				}

				_, err = s.db.QueryContext(ctx, query, params...)
				if err != nil {
					return err
				}
			}
		} else {
			query, params, err := s.db.BindNamed(
				insertTask,
				newTask{TypeID: *node.TaskTypeID, Status: "new"},
			)
			if err != nil {
				return err
			}

			_, err = s.db.QueryContext(ctx, query, params...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
