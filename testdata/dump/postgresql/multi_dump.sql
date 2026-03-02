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
(1, 'ivan.petrov@yandex.ru', 'Иван Петров', 'ivan.petrov@yandex.ru', '+7 (916) 555-12-34', 1),
(2, '(646) 555-0199', 'Emily Carter', 'emily.carter@yahoo.com', '(646) 555-0199', 2),
(3, 'lukas.schmidt', 'Lukas Schmidt', 'lukas.schmidt@web.de', '+49 30 1234 5678', 3),
(4, '01 42 68 53 01', 'Camille Bernard', 'camille.bernard@free.fr', '01 42 68 53 01', 2),
(5, 'erik.andersson@telia.se', 'Erik Andersson', 'erik.andersson@telia.se', '+46 8 123 45 67', 1);

CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +44 20 7000 1234, email editorial.office@news.co.uk. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone +33 1 55 44 33 22, email privacy.board@orange.fr. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +1 (202) 555-0141, email help.center@gmail.com. Keep this note in the exported dump.', 3);
