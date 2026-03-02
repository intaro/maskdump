-- Oracle Database dump
-- Schema: MASKDUMP_US
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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'john.miller', 'John Miller', 'j8ae554@gmail.com', '+1 (802) 235-0158', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '(673) 548-0798', 'Emily Carter', 'e526362@yahoo.com', '(673) 548-0798', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'me4e8d5@outlook.com', 'Mason Hall', 'me4e8d5@outlook.com', '498-532-0536', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'olivia.wright', 'Olivia Wright', 'o1ca5a2@proton.me', '+1 093 785 8184', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '3465230742', 'Noah Davis', 'n989ba1@icloud.com', '3465230742', 1);
CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +1 (622) 145-4131, email mfcb75d@newsroom.us. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 401-592-0192, email pf1ddf1@outlook.com. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone (691) 521-0801, email hf31544@gmail.com. Keep this note in the exported dump.', 3);
COMMIT;
