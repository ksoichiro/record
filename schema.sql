DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `task`;
DROP TABLE IF EXISTS `history`;

CREATE TABLE `user` (
  `id` int(11) primary key auto_increment,
  `name` varchar(100) not null,
  `created_at` datetime not null
) DEFAULT CHARSET=utf8;

CREATE TABLE `task` (
  `id` int(11) primary key auto_increment,
  `user_id` int(11) not null,
  `title` varchar(200) not null,
  `description` text,
  `done` tinyint(1) not null default 0,
  `type` int(11) not null,
  `amount` int(11),
  `created_at` datetime not null
) DEFAULT CHARSET=utf8;

CREATE TABLE `history` (
  `id` int(11) primary key auto_increment,
  `user_id` int(11) not null,
  `task_id` int(11) not null,
  `amount` int(11),
  `created_at` datetime not null
) DEFAULT CHARSET=utf8;
