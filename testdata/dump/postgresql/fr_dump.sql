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
(1, 'luc.martin', 'Luc Martin', 'luc.martin@orange.fr', '+33 1 42 68 53 00', 1),
(2, '01 42 68 53 01', 'Camille Bernard', 'camille.bernard@free.fr', '01 42 68 53 01', 2),
(3, 'julie.dubois@sfr.fr', 'Julie Dubois', 'julie.dubois@sfr.fr', '+33 (0)4 72 00 00 00', 3),
(4, 'nicolas.moreau', 'Nicolas Moreau', 'nicolas.moreau@laposte.net', '06 12 34 56 78', 2),
(5, '0611223344', 'Lea Petit', 'lea.petit@proton.me', '0611223344', 1);

CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +33 1 55 44 33 22, email presse@entreprise.fr. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 04 72 10 20 30, email confidentialite@orange.fr. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +33 (0)3 88 11 22 33, email support-client@free.fr. Keep this note in the exported dump.', 3);
