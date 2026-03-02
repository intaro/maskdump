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
(1, 'lukas.schmidt', 'Lukas Schmidt', 'lukas.schmidt@web.de', '+49 30 1234 5678', 1),
(2, '030 123456', 'Anna Muller', 'anna.mueller@gmx.de', '030 123456', 2),
(3, 'leonie.fischer@mail.de', 'Leonie Fischer', 'leonie.fischer@mail.de', '+49 (89) 2345 6789', 3),
(4, 'max.weber', 'Max Weber', 'max.weber@t-online.de', '040 987654', 2),
(5, '01761234567', 'Sophie Becker', 'sophie.becker@posteo.de', '01761234567', 1);

CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +49 211 4567 8910, email presse@firma.de. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 089 998877, email datenschutz@web.de. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +49 40 7654 3210, email hilfe@gmx.de. Keep this note in the exported dump.', 3);
