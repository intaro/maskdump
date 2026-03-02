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
(1, 'oliver.smith', 'Oliver Smith', 'o172b6c@btinternet.com', '+43 50 2046 0558', 1),
(2, '054 7476 2959', 'Amelia Brown', 'adc29bf@outlook.co.uk', '054 7476 2959', 2),
(3, 'h9d5504@gmail.com', 'Harry Jones', 'h9d5504@gmail.com', '+42 662 998 0000', 3),
(4, 'isla.wilson', 'Isla Wilson', 'i3f3a84@protonmail.com', '0847 956 9153', 2),
(5, '03708516272', 'George Taylor', 'geb413f@yahoo.co.uk', '03708516272', 1);
CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +44 713 195 0901, email p910e67@news.co.uk. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 079 7970 4294, email pf7cd2f@outlook.co.uk. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +48 125 790 0802, email h288682@gmail.com. Keep this note in the exported dump.', 3);
