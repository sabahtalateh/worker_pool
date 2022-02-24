-- +goose Up
-- +goose StatementBegin

create table jobs
(
    id           serial primary key,
    template_id  int                                    not null,
    submitted_at timestamp with time zone default now() not null,

    constraint deployment_job_template_fk foreign key (template_id) references job_templates (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table jobs;

-- +goose StatementEnd
