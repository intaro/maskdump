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
(1, 'ivan.petrov', 'Иван Петров', 'ivan.petrov@yandex.ru', '+7 (916) 555-12-34', 1),
(2, '8 912 444 55 66', 'Анна Смирнова', 'anna.smirnova@mail.ru', '8 912 444 55 66', 2),
(3, 'sergey-volkov', 'Сергей Волков', 'sergey.volkov@bk.ru', '7 495 123 45 67', 3),
(4, 'olga.romanova@rambler.ru', 'Ольга Романова', 'olga.romanova@rambler.ru', '+7 812 600 77 88', 2),
(5, '79165550011', 'Елена Соколова', 'elena.sokolova@list.ru', '79165550011', 1);

CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +7 (495) 777-11-22, email press-office@company.ru. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 8 800 555 35 35, email privacy-team@yandex.ru. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone 7 812 320 10 10, email support-center@mail.ru. Keep this note in the exported dump.', 3);
