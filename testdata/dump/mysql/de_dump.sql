-- MySQL dump 10.13  Distrib 8.0.36, for Linux (x86_64)
--
-- Host: localhost    Database: maskdump_fixture_de
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
(1,'lukas.schmidt','Lukas Schmidt','lukas.schmidt@web.de','+49 30 1234 5678',1),
(2,'030 123456','Anna Muller','anna.mueller@gmx.de','030 123456',2),
(3,'leonie.fischer@mail.de','Leonie Fischer','leonie.fischer@mail.de','+49 (89) 2345 6789',3),
(4,'max.weber','Max Weber','max.weber@t-online.de','040 987654',2),
(5,'01761234567','Sophie Becker','sophie.becker@posteo.de','01761234567',1);

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
(1,'welcome-playbook','Welcome Playbook','Escalation contact 1: phone +49 211 4567 8910, email presse@firma.de. Keep this note in the exported dump.',1),
(2,'privacy-checklist','Privacy Checklist','Escalation contact 2: phone 089 998877, email datenschutz@web.de. Keep this note in the exported dump.',2),
(3,'support-handbook','Support Handbook','Escalation contact 3: phone +49 40 7654 3210, email hilfe@gmx.de. Keep this note in the exported dump.',3);
