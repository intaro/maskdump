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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'lukas.schmidt', 'Lukas Schmidt', 'lba36a4@web.de', '+45 50 2536 5878', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '050 123456', 'Anna Muller', 'a6f9158@gmx.de', '050 123456', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'l032979@mail.de', 'Leonie Fischer', 'l032979@mail.de', '+43 (59) 2540 6389', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'max.weber', 'Max Weber', 'm6aaafc@t-online.de', '017 980644', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '01664836507', 'Sophie Becker', 's3e7368@posteo.de', '01664836507', 1);
CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +43 111 3587 2910, email pb1187e@firma.de. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 053 915877, email d28a13b@web.de. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +45 40 2252 3510, email hf9af31@gmx.de. Keep this note in the exported dump.', 3);
COMMIT;
