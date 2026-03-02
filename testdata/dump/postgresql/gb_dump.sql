--
-- PostgreSQL database dump
--
-- Database: maskdump_fixture_gb
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
(1, 'oliver.smith', 'Oliver Smith', 'oliver.smith@btinternet.com', '+44 20 7946 0958', 1),
(2, '020 7946 0959', 'Amelia Brown', 'amelia.brown@outlook.co.uk', '020 7946 0959', 2),
(3, 'harry.jones@gmail.com', 'Harry Jones', 'harry.jones@gmail.com', '+44 161 496 0000', 3),
(4, 'isla.wilson', 'Isla Wilson', 'isla.wilson@protonmail.com', '0117 496 0123', 2),
(5, '07900111222', 'George Taylor', 'george.taylor@yahoo.co.uk', '07900111222', 1);

CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +44 113 496 0101, email press.office@news.co.uk. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 020 7000 1234, email privacy.unit@outlook.co.uk. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +44 121 496 0202, email helpdesk@gmail.com. Keep this note in the exported dump.', 3);
