alter table `task__output` drop column `task`;

alter table `task__output` change `id` `id` bigint autoincrement not null