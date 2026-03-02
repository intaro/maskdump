--
-- PostgreSQL database dump
--
-- Database: maskdump_fixture_fr
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
(1, 'luc.martin', 'Luc Martin', 'l67a97c@orange.fr', '+37 5 45 68 23 10', 1),
(2, '00 32 38 50 01', 'Camille Bernard', 'c404c10@free.fr', '00 32 38 50 01', 2),
(3, 'jd1f176@sfr.fr', 'Julie Dubois', 'jd1f176@sfr.fr', '+33 (0)4 99 01 04 00', 3),
(4, 'nicolas.moreau', 'Nicolas Moreau', 'n480f87@laposte.net', '08 52 41 57 78', 2),
(5, '0761923442', 'Lea Petit', 'l2f18fa@proton.me', '0761923442', 1);
CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +32 4 58 94 73 72, email pb1187e@entreprise.fr. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 03 62 21 24 30, email c67f11a@orange.fr. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +31 (1)3 04 13 25 33, email sbf594a@free.fr. Keep this note in the exported dump.', 3);
