update `project__template` set `app` = 'ansible' where `app` = '';

--Already set in 2.9.46
--alter table `project__template` change `app` `app` varchar(50) not null;
