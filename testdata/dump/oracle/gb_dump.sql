-- Oracle Database dump
-- Schema: MASKDUMP_GB
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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'oliver.smith', 'Oliver Smith', 'oliver.smith@btinternet.com', '+44 20 7946 0958', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '020 7946 0959', 'Amelia Brown', 'amelia.brown@outlook.co.uk', '020 7946 0959', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'harry.jones@gmail.com', 'Harry Jones', 'harry.jones@gmail.com', '+44 161 496 0000', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'isla.wilson', 'Isla Wilson', 'isla.wilson@protonmail.com', '0117 496 0123', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '07900111222', 'George Taylor', 'george.taylor@yahoo.co.uk', '07900111222', 1);

CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +44 113 496 0101, email press.office@news.co.uk. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 020 7000 1234, email privacy.unit@outlook.co.uk. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +44 121 496 0202, email helpdesk@gmail.com. Keep this note in the exported dump.', 3);
COMMIT;
