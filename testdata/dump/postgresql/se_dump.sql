--
-- PostgreSQL database dump
--
-- Database: maskdump_fixture_se
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
(1, 'erik.andersson', 'Erik Andersson', 'erik.andersson@telia.se', '+46 8 123 45 67', 1),
(2, '08-123 45 68', 'Anna Johansson', 'anna.johansson@outlook.com', '08-123 45 68', 2),
(3, 'elsa.nilsson@gmail.com', 'Elsa Nilsson', 'elsa.nilsson@gmail.com', '+46 (0)31-123 456', 3),
(4, 'oscar.lindberg', 'Oscar Lindberg', 'oscar.lindberg@bahnhof.se', '031-765 432', 2),
(5, '0701234567', 'Maja Karlsson', 'maja.karlsson@icloud.com', '0701234567', 1);

CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +46 31 701 23 45, email press@bolag.se. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 08-555 12 34, email privacy.office@telia.se. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +46 (0)90-123 456, email kundservice@bahnhof.se. Keep this note in the exported dump.', 3);
