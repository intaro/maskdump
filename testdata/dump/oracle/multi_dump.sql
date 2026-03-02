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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'ivan.petrov@yandex.ru', 'Иван Петров', 'ivan.petrov@yandex.ru', '+7 (916) 555-12-34', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '(646) 555-0199', 'Emily Carter', 'emily.carter@yahoo.com', '(646) 555-0199', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'lukas.schmidt', 'Lukas Schmidt', 'lukas.schmidt@web.de', '+49 30 1234 5678', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, '01 42 68 53 01', 'Camille Bernard', 'camille.bernard@free.fr', '01 42 68 53 01', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, 'erik.andersson@telia.se', 'Erik Andersson', 'erik.andersson@telia.se', '+46 8 123 45 67', 1);

CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +44 20 7000 1234, email editorial.office@news.co.uk. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone +33 1 55 44 33 22, email privacy.board@orange.fr. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +1 (202) 555-0141, email help.center@gmail.com. Keep this note in the exported dump.', 3);
COMMIT;
