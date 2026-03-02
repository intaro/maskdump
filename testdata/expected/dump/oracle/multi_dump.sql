-- Oracle Database dump
-- Schema: MASKDUMP_MULTI
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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'i690245@yandex.ru', 'Иван Петров', 'i690245@yandex.ru', '+7 (506) 405-52-34', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '(673) 548-0798', 'Emily Carter', 'e526362@yahoo.com', '(673) 548-0798', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'lukas.schmidt', 'Lukas Schmidt', 'lba36a4@web.de', '+45 50 2536 5878', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, '00 32 38 50 01', 'Camille Bernard', 'c404c10@free.fr', '00 32 38 50 01', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, 'ed763e6@telia.se', 'Erik Andersson', 'ed763e6@telia.se', '+48 7 135 40 65', 1);
CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +44 50 4504 1634, email e26860b@news.co.uk. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone +32 4 58 94 73 72, email p514224@orange.fr. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +1 (622) 145-4131, email hf31544@gmail.com. Keep this note in the exported dump.', 3);
COMMIT;
