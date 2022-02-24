-- +goose Up
-- +goose StatementBegin

create table task_types
(
    id         serial primary key,

    name       varchar                                not null,

    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null
);

insert into task_types (name)
values ('wait_30_secs'),
       ('wait_1_min'),
       ('build'),
       ('clone_config'),
       ('render_config'),
       ('deploy')
;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table task_types;

-- +goose StatementEnd
