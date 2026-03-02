-- MySQL dump 10.13  Distrib 8.0.36, for Linux (x86_64)
--
-- Host: localhost    Database: maskdump_fixture_fr
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
(1,'luc.martin','Luc Martin','l67a97c@orange.fr','+37 5 45 68 23 10',1),
(2,'00 32 38 50 01','Camille Bernard','c404c10@free.fr','00 32 38 50 01',2),
(3,'jd1f176@sfr.fr','Julie Dubois','jd1f176@sfr.fr','+33 (0)4 99 01 04 00',3),
(4,'nicolas.moreau','Nicolas Moreau','n480f87@laposte.net','08 52 41 57 78',2),
(5,'0761923442','Lea Petit','l2f18fa@proton.me','0761923442',1);
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
(1,'welcome-playbook','Welcome Playbook','Escalation contact 1: phone +32 4 58 94 73 72, email pb1187e@entreprise.fr. Keep this note in the exported dump.',1),
(2,'privacy-checklist','Privacy Checklist','Escalation contact 2: phone 03 62 21 24 30, email c67f11a@orange.fr. Keep this note in the exported dump.',2),
(3,'support-handbook','Support Handbook','Escalation contact 3: phone +31 (1)3 04 13 25 33, email sbf594a@free.fr. Keep this note in the exported dump.',3);
