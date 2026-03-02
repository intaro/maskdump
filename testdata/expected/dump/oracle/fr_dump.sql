-- Oracle Database dump
-- Schema: MASKDUMP_FR
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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'luc.martin', 'Luc Martin', 'l67a97c@orange.fr', '+37 5 45 68 23 10', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '00 32 38 50 01', 'Camille Bernard', 'c404c10@free.fr', '00 32 38 50 01', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'jd1f176@sfr.fr', 'Julie Dubois', 'jd1f176@sfr.fr', '+33 (0)4 99 01 04 00', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'nicolas.moreau', 'Nicolas Moreau', 'n480f87@laposte.net', '08 52 41 57 78', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '0761923442', 'Lea Petit', 'l2f18fa@proton.me', '0761923442', 1);
CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +32 4 58 94 73 72, email pb1187e@entreprise.fr. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 03 62 21 24 30, email c67f11a@orange.fr. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +31 (1)3 04 13 25 33, email sbf594a@free.fr. Keep this note in the exported dump.', 3);
COMMIT;
