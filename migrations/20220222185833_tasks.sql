-- +goose Up
-- +goose StatementBegin

create type task_status as enum ('new', 'progressing', 'success', 'failure');

create table tasks
(
    id           serial primary key,

    type_id      int                                    not null,

    status       task_status                            not null default 'new',

    created_at   timestamp with time zone default now() not null,
    updated_at   timestamp with time zone default now() not null,
    requested_at timestamp with time zone default now() not null,

    constraint fk_type_id foreign key (type_id) references task_types (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table tasks;
drop type task_status;

-- +goose StatementEnd
