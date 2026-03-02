--
-- PostgreSQL database dump
--
-- Database: maskdump_fixture_de
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
(1, 'lukas.schmidt', 'Lukas Schmidt', 'lba36a4@web.de', '+45 50 2536 5878', 1),
(2, '050 123456', 'Anna Muller', 'a6f9158@gmx.de', '050 123456', 2),
(3, 'l032979@mail.de', 'Leonie Fischer', 'l032979@mail.de', '+43 (59) 2540 6389', 3),
(4, 'max.weber', 'Max Weber', 'm6aaafc@t-online.de', '017 980644', 2),
(5, '01664836507', 'Sophie Becker', 's3e7368@posteo.de', '01664836507', 1);
CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +43 111 3587 2910, email pb1187e@firma.de. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 053 915877, email d28a13b@web.de. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +45 40 2252 3510, email hf9af31@gmx.de. Keep this note in the exported dump.', 3);
