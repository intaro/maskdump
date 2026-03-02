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
(1, 'erik.andersson', 'Erik Andersson', 'ed763e6@telia.se', '+48 7 135 40 65', 1),
(2, '09-921 15 48', 'Anna Johansson', 'a098626@outlook.com', '09-921 15 48', 2),
(3, 'eaafd9e@gmail.com', 'Elsa Nilsson', 'eaafd9e@gmail.com', '+46 (0)37-224 559', 3),
(4, 'oscar.lindberg', 'Oscar Lindberg', 'o8b7f6b@bahnhof.se', '070-715 472', 2),
(5, '0911374868', 'Maja Karlsson', 'mf2f5a6@icloud.com', '0911374868', 1);
CREATE TABLE public.tst_posts (
    id bigint PRIMARY KEY,
    code varchar(128) NOT NULL,
    title varchar(255) NOT NULL,
    detail text NOT NULL,
    user_id bigint NOT NULL
);
INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES
(1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +45 61 641 33 45, email p5f0038@bolag.se. Keep this note in the exported dump.', 1),
(2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 05-756 22 14, email p1989f1@telia.se. Keep this note in the exported dump.', 2),
(3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +46 (0)99-321 553, email k28620f@bahnhof.se. Keep this note in the exported dump.', 3);
