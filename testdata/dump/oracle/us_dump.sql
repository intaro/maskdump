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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'john.miller', 'John Miller', 'john.miller@gmail.com', '+1 (212) 555-0188', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '(646) 555-0199', 'Emily Carter', 'emily.carter@yahoo.com', '(646) 555-0199', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'mason.hall@outlook.com', 'Mason Hall', 'mason.hall@outlook.com', '415-555-0132', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'olivia.wright', 'Olivia Wright', 'olivia.wright@proton.me', '+1 503 555 0114', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '3125550147', 'Noah Davis', 'noah.davis@icloud.com', '3125550147', 1);

CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +1 (202) 555-0141, email media.desk@newsroom.us. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 415-555-0198, email privacy-team@outlook.com. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone (646) 555-0102, email help.center@gmail.com. Keep this note in the exported dump.', 3);
COMMIT;
