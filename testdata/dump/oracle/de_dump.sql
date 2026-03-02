-- Oracle Database dump
-- Schema: MASKDUMP_DE
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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'lukas.schmidt', 'Lukas Schmidt', 'lukas.schmidt@web.de', '+49 30 1234 5678', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '030 123456', 'Anna Muller', 'anna.mueller@gmx.de', '030 123456', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'leonie.fischer@mail.de', 'Leonie Fischer', 'leonie.fischer@mail.de', '+49 (89) 2345 6789', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'max.weber', 'Max Weber', 'max.weber@t-online.de', '040 987654', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '01761234567', 'Sophie Becker', 'sophie.becker@posteo.de', '01761234567', 1);

CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +49 211 4567 8910, email presse@firma.de. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 089 998877, email datenschutz@web.de. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +49 40 7654 3210, email hilfe@gmx.de. Keep this note in the exported dump.', 3);
COMMIT;
