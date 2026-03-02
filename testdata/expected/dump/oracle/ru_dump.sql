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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'ivan.petrov', 'Иван Петров', 'i690245@yandex.ru', '+7 (506) 405-52-34', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '8 212 974 15 06', 'Анна Смирнова', 'ad54e29@mail.ru', '8 212 974 15 06', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'sergey-volkov', 'Сергей Волков', 's608e39@bk.ru', '7 015 893 35 37', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'o87cc5e@rambler.ru', 'Ольга Романова', 'o87cc5e@rambler.ru', '+7 012 200 17 88', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '72962859031', 'Елена Соколова', 'e18ddbb@list.ru', '72962859031', 1);
CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +7 (375) 297-51-52, email pcd9ca7@company.ru. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 8 560 205 45 95, email pf1ddf1@yandex.ru. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone 7 062 480 30 70, email sef0d2d@mail.ru. Keep this note in the exported dump.', 3);
COMMIT;
