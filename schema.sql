DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `tasks`;
DROP TABLE IF EXISTS `histories`;

CREATE TABLE `users` (
  `id` int(11) primary key auto_increment,
  `name` varchar(100) not null,
  `password` varchar(100) not null,
  `created_at` datetime not null
) DEFAULT CHARSET=utf8;
ALTER TABLE `users` ADD UNIQUE `uq_users` (`name`);

CREATE TABLE `tasks` (
  `id` int(11) primary key auto_increment,
  `user_id` int(11) not null,
  `title` varchar(200) not null,
  `description` text,
  `done` tinyint(1) not null default 0,
  `type` int(11) not null,
  `amount` int(11),
  `created_at` datetime not null
) DEFAULT CHARSET=utf8;

CREATE TABLE `histories` (
  `id` int(11) primary key auto_increment,
  `user_id` int(11) not null,
  `task_id` int(11) not null,
  `amount` int(11),
  `created_at` datetime not null
) DEFAULT CHARSET=utf8;
