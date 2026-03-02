-- Oracle Database dump
-- Schema: MASKDUMP_RU
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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'ivan.petrov', 'Иван Петров', 'ivan.petrov@yandex.ru', '+7 (916) 555-12-34', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '8 912 444 55 66', 'Анна Смирнова', 'anna.smirnova@mail.ru', '8 912 444 55 66', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'sergey-volkov', 'Сергей Волков', 'sergey.volkov@bk.ru', '7 495 123 45 67', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'olga.romanova@rambler.ru', 'Ольга Романова', 'olga.romanova@rambler.ru', '+7 812 600 77 88', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '79165550011', 'Елена Соколова', 'elena.sokolova@list.ru', '79165550011', 1);

CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +7 (495) 777-11-22, email press-office@company.ru. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 8 800 555 35 35, email privacy-team@yandex.ru. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone 7 812 320 10 10, email support-center@mail.ru. Keep this note in the exported dump.', 3);
COMMIT;
