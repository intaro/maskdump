--
-- PostgreSQL database dump
--
-- Database: maskdump_fixture_ru
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
(1, 'ivan.petrov', 'Иван Петров', 'i690245@yandex.ru', '+7 (506) 405-52-34', 1),
(2, '8 212 974 15 06', 'Анна Смирнова', 'ad54e29@mail.ru', '8 212 974 15 06', 2),
(3, 'sergey-volkov', 'Сергей Волков', 's608e39@bk.ru', '7 015 893 35 37', 3),
(4, 'o87cc5e@rambler.ru', 'Ольга Романова', 'o87cc5e@rambler.ru', '+7 012 200 17 88', 2),
(5, '72962859031', 'Елена Соколова', 'e18ddbb@list.ru', '72962859031', 1);
CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +7 (375) 297-51-52, email pcd9ca7@company.ru. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 8 560 205 45 95, email pf1ddf1@yandex.ru. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone 7 062 480 30 70, email sef0d2d@mail.ru. Keep this note in the exported dump.', 3);
