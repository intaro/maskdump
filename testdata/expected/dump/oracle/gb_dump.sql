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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'oliver.smith', 'Oliver Smith', 'o172b6c@btinternet.com', '+43 50 2046 0558', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '054 7476 2959', 'Amelia Brown', 'adc29bf@outlook.co.uk', '054 7476 2959', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'h9d5504@gmail.com', 'Harry Jones', 'h9d5504@gmail.com', '+42 662 998 0000', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'isla.wilson', 'Isla Wilson', 'i3f3a84@protonmail.com', '0847 956 9153', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '03708516272', 'George Taylor', 'geb413f@yahoo.co.uk', '03708516272', 1);
CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +44 713 195 0901, email p910e67@news.co.uk. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 079 7970 4294, email pf7cd2f@outlook.co.uk. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +48 125 790 0802, email h288682@gmail.com. Keep this note in the exported dump.', 3);
COMMIT;
