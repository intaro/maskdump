-- Oracle Database dump
-- Schema: MASKDUMP_SE
SET DEFINE OFF;
BEGIN EXECUTE IMMEDIATE 'DROP TABLE tst_groups'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE tst_users'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE tst_posts'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
CREATE TABLE tst_groups (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(64 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL
);
INSERT INTO tst_groups (id, code, title) VALUES (1, 'admins', 'Administrators');
INSERT INTO tst_groups (id, code, title) VALUES (2, 'editors', 'Editorial Team');
INSERT INTO tst_groups (id, code, title) VALUES (3, 'support', 'Customer Success');
CREATE TABLE tst_users (
    id NUMBER(19) PRIMARY KEY,
    login VARCHAR2(255 CHAR) NOT NULL,
    name VARCHAR2(255 CHAR) NOT NULL,
    email VARCHAR2(255 CHAR) NOT NULL,
    phone VARCHAR2(255 CHAR) NOT NULL,
    group_id NUMBER(19) NOT NULL
);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'erik.andersson', 'Erik Andersson', 'ed763e6@telia.se', '+48 7 135 40 65', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '09-921 15 48', 'Anna Johansson', 'a098626@outlook.com', '09-921 15 48', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'eaafd9e@gmail.com', 'Elsa Nilsson', 'eaafd9e@gmail.com', '+46 (0)37-224 559', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'oscar.lindberg', 'Oscar Lindberg', 'o8b7f6b@bahnhof.se', '070-715 472', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '0911374868', 'Maja Karlsson', 'mf2f5a6@icloud.com', '0911374868', 1);
CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +45 61 641 33 45, email p5f0038@bolag.se. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 05-756 22 14, email p1989f1@telia.se. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +46 (0)99-321 553, email k28620f@bahnhof.se. Keep this note in the exported dump.', 3);
COMMIT;
