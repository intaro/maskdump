-- Oracle Database dump
-- Schema: MASKDUMP_SE
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
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (1, 'erik.andersson', 'Erik Andersson', 'erik.andersson@telia.se', '+46 8 123 45 67', 1);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (2, '08-123 45 68', 'Anna Johansson', 'anna.johansson@outlook.com', '08-123 45 68', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (3, 'elsa.nilsson@gmail.com', 'Elsa Nilsson', 'elsa.nilsson@gmail.com', '+46 (0)31-123 456', 3);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (4, 'oscar.lindberg', 'Oscar Lindberg', 'oscar.lindberg@bahnhof.se', '031-765 432', 2);
INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (5, '0701234567', 'Maja Karlsson', 'maja.karlsson@icloud.com', '0701234567', 1);

CREATE TABLE tst_posts (
    id NUMBER(19) PRIMARY KEY,
    code VARCHAR2(128 CHAR) NOT NULL,
    title VARCHAR2(255 CHAR) NOT NULL,
    detail CLOB NOT NULL,
    user_id NUMBER(19) NOT NULL
);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (1, 'welcome-playbook', 'Welcome Playbook', 'Escalation contact 1: phone +46 31 701 23 45, email press@bolag.se. Keep this note in the exported dump.', 1);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (2, 'privacy-checklist', 'Privacy Checklist', 'Escalation contact 2: phone 08-555 12 34, email privacy.office@telia.se. Keep this note in the exported dump.', 2);
INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (3, 'support-handbook', 'Support Handbook', 'Escalation contact 3: phone +46 (0)90-123 456, email kundservice@bahnhof.se. Keep this note in the exported dump.', 3);
COMMIT;
