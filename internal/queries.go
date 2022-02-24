package internal

const selectTaskId = `
select id
from tasks
where status = 'new'
limit 1 for update skip locked`

const selectJobTemplate = `select id, name from job_templates where name = $1`

const selectJobTemplateNodes = `select id,
       root,
       name,
       template_id,
       type,
       apply_to,
       composition,
       g1,
       g2,
       task_type_id
from job_template_nodes
where type = 'task' and template_id = $1`

const insertTask = `insert into tasks (type_id, status, service_id) values (:type_id, :status, :service_id)`
