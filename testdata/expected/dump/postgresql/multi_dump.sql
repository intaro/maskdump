--
-- PostgreSQL database dump
--
-- Database: maskdump_fixture_multi
SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
CREATE TABLE public.tst_groups (
    id bigint PRIMARY KEY,
    code varchar(64) NOT NULL,
    title varchar(255) NOT NULL
);
INSERT INTO public.tst_groups (id, code, title) VALUES
(1, 'admins', 'Administrators'),
(2, 'editors', 'Editorial Team'),
(3, 'support', 'Customer Success');
CREATE TABLE public.tst_users (
    id bigint PRIMARY KEY,
    login varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    phone varchar(255) NOT NULL,
    group_id bigint NOT NULL
);
INSERT INTO public.tst_users (id, login, name, email, phone, group_id) VALUES
(1, 'i690245@yandex.ru', 'Иван Петров', 'i690245@yandex.ru', '+7 (506) 405-52-34', 1),
(2, '(673) 548-0798', 'Emily Carter', 'e526362@yahoo.com', '(673) 548-0798', 2),
(3, 'lukas.schmidt', 'Lukas Schmidt', 'lba36a4@web.de', '+45 50 2536 5878', 3),
(4, '00 32 38 50 01', 'Camille Bernard', 'c404c10@free.fr', '00 32 38 50 01', 2),
(5, 'ed763e6@telia.se', 'Erik Andersson', 'ed763e6@telia.se', '+48 7 135 40 65', 1);
CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +44 50 4504 1634, email e26860b@news.co.uk. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone +32 4 58 94 73 72, email p514224@orange.fr. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +1 (622) 145-4131, email hf31544@gmail.com. Keep this note in the exported dump.', 3);
