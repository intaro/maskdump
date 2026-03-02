--
-- PostgreSQL database dump
--
-- Database: maskdump_fixture_us
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
(1, 'john.miller', 'John Miller', 'john.miller@gmail.com', '+1 (212) 555-0188', 1),
(2, '(646) 555-0199', 'Emily Carter', 'emily.carter@yahoo.com', '(646) 555-0199', 2),
(3, 'mason.hall@outlook.com', 'Mason Hall', 'mason.hall@outlook.com', '415-555-0132', 3),
(4, 'olivia.wright', 'Olivia Wright', 'olivia.wright@proton.me', '+1 503 555 0114', 2),
(5, '3125550147', 'Noah Davis', 'noah.davis@icloud.com', '3125550147', 1);

CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +1 (202) 555-0141, email media.desk@newsroom.us. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 415-555-0198, email privacy-team@outlook.com. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone (646) 555-0102, email help.center@gmail.com. Keep this note in the exported dump.', 3);
