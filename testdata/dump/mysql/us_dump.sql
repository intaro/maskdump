-- MySQL dump 10.13  Distrib 8.0.36, for Linux (x86_64)
--
-- Host: localhost    Database: maskdump_fixture_us
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
(1,'john.miller','John Miller','john.miller@gmail.com','+1 (212) 555-0188',1),
(2,'(646) 555-0199','Emily Carter','emily.carter@yahoo.com','(646) 555-0199',2),
(3,'mason.hall@outlook.com','Mason Hall','mason.hall@outlook.com','415-555-0132',3),
(4,'olivia.wright','Olivia Wright','olivia.wright@proton.me','+1 503 555 0114',2),
(5,'3125550147','Noah Davis','noah.davis@icloud.com','3125550147',1);

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
(1,'welcome-playbook','Welcome Playbook','Escalation contact 1: phone +1 (202) 555-0141, email media.desk@newsroom.us. Keep this note in the exported dump.',1),
(2,'privacy-checklist','Privacy Checklist','Escalation contact 2: phone 415-555-0198, email privacy-team@outlook.com. Keep this note in the exported dump.',2),
(3,'support-handbook','Support Handbook','Escalation contact 3: phone (646) 555-0102, email help.center@gmail.com. Keep this note in the exported dump.',3);
