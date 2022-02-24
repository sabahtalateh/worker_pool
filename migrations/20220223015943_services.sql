-- +goose Up
-- +goose StatementBegin

create table services
(
    id   serial primary key,
    name varchar not null
);

insert into services (name)
values ('table-booking');

insert into services (name)
values ('cmrs');

insert into services (name)
values ('vendor-onboarding');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table services;

-- +goose StatementEnd
