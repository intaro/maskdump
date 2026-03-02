-- MySQL dump 10.13  Distrib 8.0.36, for Linux (x86_64)
--
-- Host: localhost    Database: maskdump_fixture_ru
-- ------------------------------------------------------
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8mb4 */;
DROP TABLE IF EXISTS `tst_groups`;
CREATE TABLE `tst_groups` (
  `id` bigint NOT NULL,
  `code` varchar(64) NOT NULL,
  `title` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT INTO `tst_groups` VALUES
(1,'admins','Administrators'),
(2,'editors','Editorial Team'),
(3,'support','Customer Success');
DROP TABLE IF EXISTS `tst_users`;
CREATE TABLE `tst_users` (
  `id` bigint NOT NULL,
  `login` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `phone` varchar(255) NOT NULL,
  `group_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT INTO `tst_users` VALUES
(1,'ivan.petrov','Иван Петров','i690245@yandex.ru','+7 (506) 405-52-34',1),
(2,'8 212 974 15 06','Анна Смирнова','ad54e29@mail.ru','8 212 974 15 06',2),
(3,'sergey-volkov','Сергей Волков','s608e39@bk.ru','7 015 893 35 37',3),
(4,'o87cc5e@rambler.ru','Ольга Романова','o87cc5e@rambler.ru','+7 012 200 17 88',2),
(5,'72962859031','Елена Соколова','e18ddbb@list.ru','72962859031',1);
DROP TABLE IF EXISTS `tst_posts`;
CREATE TABLE `tst_posts` (
  `id` bigint NOT NULL,
  `code` varchar(128) NOT NULL,
  `title` varchar(255) NOT NULL,
  `detail` text NOT NULL,
  `user_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT INTO `tst_posts` VALUES
(1,'welcome-playbook','Welcome Playbook','Escalation contact 1: phone +7 (375) 297-51-52, email pcd9ca7@company.ru. Keep this note in the exported dump.',1),
(2,'privacy-checklist','Privacy Checklist','Escalation contact 2: phone 8 560 205 45 95, email pf1ddf1@yandex.ru. Keep this note in the exported dump.',2),
(3,'support-handbook','Support Handbook','Escalation contact 3: phone 7 062 480 30 70, email sef0d2d@mail.ru. Keep this note in the exported dump.',3);
