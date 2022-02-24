-- +goose Up
-- +goose StatementBegin

create table job_templates
(
    id   serial primary key,
    name varchar not null
);

create type job_template_node_type as enum ('group', 'task');
create type job_template_node_composition as enum ('parallel', 'consequence');
create type job_template_target as enum ('single_service', 'all_services');

create table job_template_nodes
(
    id           serial primary key,
    root         bool                   not null default false,
    name         varchar unique         not null,
    template_id  int                    not null,
    type         job_template_node_type not null,
    apply_to     job_template_target    not null default 'single_service',
    composition  job_template_node_composition, -- in case of type == group
    g1           int,                           -- in case of type == group
    g2           int,                           -- in case of type == group
    task_type_id int,                           -- in case of type == task

    constraint job_template_id foreign key (template_id) references job_templates (id),
    constraint job_g1_node foreign key (g1) references job_template_nodes (id),
    constraint job_g2_node foreign key (g2) references job_template_nodes (id),
    constraint node_task_type_id foreign key (task_type_id) references task_types (id)
);

insert into job_templates (name)
values ('standard deployment');

insert into job_template_nodes (template_id, name, type, composition, g1, g2, task_type_id)
values ((select id from job_templates where name = 'standard deployment'), 'build', 'task', null, null, null,
        (select id from task_types where name = 'build'));

insert into job_template_nodes (template_id, name, type, composition, g1, g2, task_type_id)
values ((select id from job_templates where name = 'standard deployment'), 'clone config', 'task', null, null, null,
        (select id from task_types where name = 'clone_config'));

insert into job_template_nodes (template_id, name, type, composition, g1, g2, task_type_id)
values ((select id from job_templates where name = 'standard deployment'), 'render config', 'task', null, null, null,
        (select id from task_types where name = 'render_config'));

insert into job_template_nodes (template_id, name, type, composition, g1, g2)
values ((select id from job_templates where name = 'standard deployment'), 'clone & render', 'group', 'consequence',
        (select id from job_template_nodes where name = 'clone config'),
        (select id from job_template_nodes where name = 'render config'));

insert into job_template_nodes (template_id, name, type, composition, g1, g2)
values ((select id from job_templates where name = 'standard deployment'), 'build & clone & render', 'group',
        'parallel',
        (select id from job_template_nodes where name = 'build'),
        (select id from job_template_nodes where name = 'clone & render'));

insert into job_template_nodes (template_id, name, type, task_type_id)
values ((select id from job_templates where name = 'standard deployment'), 'deploy', 'task', (select id from task_types where name = 'deploy'));

insert into job_template_nodes (root, template_id, name, type, composition, g1, g2)
values (true, (select id from job_templates where name = 'standard deployment'), 'full standard deployment', 'group',
        'consequence',
        (select id from job_template_nodes where name = 'build & clone & render'),
        (select id from job_template_nodes where name = 'deploy'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table job_template_nodes;
drop table job_templates;
drop type job_template_node_type;
drop type job_template_node_composition;
drop type job_template_target;

-- +goose StatementEnd
