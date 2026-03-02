USE [maskdump_fixture_ru]
GO
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[tst_groups](
    [id] bigint NOT NULL,
    [code] nvarchar(64) NOT NULL,
    [title] nvarchar(255) NOT NULL,
    CONSTRAINT [PK_tst_groups] PRIMARY KEY CLUSTERED ([id] ASC)
)
GO
INSERT INTO [dbo].[tst_groups] ([id], [code], [title]) VALUES (1, N'admins', N'Administrators')
INSERT INTO [dbo].[tst_groups] ([id], [code], [title]) VALUES (2, N'editors', N'Editorial Team')
INSERT INTO [dbo].[tst_groups] ([id], [code], [title]) VALUES (3, N'support', N'Customer Success')
GO

CREATE TABLE [dbo].[tst_users](
    [id] bigint NOT NULL,
    [login] nvarchar(255) NOT NULL,
    [name] nvarchar(255) NOT NULL,
    [email] nvarchar(255) NOT NULL,
    [phone] nvarchar(255) NOT NULL,
    [group_id] bigint NOT NULL,
    CONSTRAINT [PK_tst_users] PRIMARY KEY CLUSTERED ([id] ASC)
)
GO
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (1, N'ivan.petrov', N'Иван Петров', N'ivan.petrov@yandex.ru', N'+7 (916) 555-12-34', 1)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (2, N'8 912 444 55 66', N'Анна Смирнова', N'anna.smirnova@mail.ru', N'8 912 444 55 66', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (3, N'sergey-volkov', N'Сергей Волков', N'sergey.volkov@bk.ru', N'7 495 123 45 67', 3)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (4, N'olga.romanova@rambler.ru', N'Ольга Романова', N'olga.romanova@rambler.ru', N'+7 812 600 77 88', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (5, N'79165550011', N'Елена Соколова', N'elena.sokolova@list.ru', N'79165550011', 1)
GO

CREATE TABLE [dbo].[tst_posts](
    [id] bigint NOT NULL,
    [code] nvarchar(128) NOT NULL,
    [title] nvarchar(255) NOT NULL,
    [detail] nvarchar(max) NOT NULL,
    [user_id] bigint NOT NULL,
    CONSTRAINT [PK_tst_posts] PRIMARY KEY CLUSTERED ([id] ASC)
)
GO
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (1, N'welcome-playbook', N'Welcome Playbook', N'Escalation contact 1: phone +7 (495) 777-11-22, email press-office@company.ru. Keep this note in the exported dump.', 1)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (2, N'privacy-checklist', N'Privacy Checklist', N'Escalation contact 2: phone 8 800 555 35 35, email privacy-team@yandex.ru. Keep this note in the exported dump.', 2)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (3, N'support-handbook', N'Support Handbook', N'Escalation contact 3: phone 7 812 320 10 10, email support-center@mail.ru. Keep this note in the exported dump.', 3)
GO
