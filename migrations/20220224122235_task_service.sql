-- +goose Up
-- +goose StatementBegin

alter table tasks
    add service_id int;
alter table tasks
    add constraint task_service_id foreign key (service_id) references services (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table tasks drop column service_id;

-- +goose StatementEnd
