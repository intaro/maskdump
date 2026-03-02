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
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (1, N'ivan.petrov', N'Иван Петров', N'i690245@yandex.ru', N'+7 (506) 405-52-34', 1)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (2, N'8 212 974 15 06', N'Анна Смирнова', N'ad54e29@mail.ru', N'8 212 974 15 06', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (3, N'sergey-volkov', N'Сергей Волков', N's608e39@bk.ru', N'7 015 893 35 37', 3)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (4, N'o87cc5e@rambler.ru', N'Ольга Романова', N'o87cc5e@rambler.ru', N'+7 012 200 17 88', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (5, N'72962859031', N'Елена Соколова', N'e18ddbb@list.ru', N'72962859031', 1)
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
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (1, N'welcome-playbook', N'Welcome Playbook', N'Escalation contact 1: phone +7 (375) 297-51-52, email pcd9ca7@company.ru. Keep this note in the exported dump.', 1)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (2, N'privacy-checklist', N'Privacy Checklist', N'Escalation contact 2: phone 8 560 205 45 95, email pf1ddf1@yandex.ru. Keep this note in the exported dump.', 2)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (3, N'support-handbook', N'Support Handbook', N'Escalation contact 3: phone 7 062 480 30 70, email sef0d2d@mail.ru. Keep this note in the exported dump.', 3)
GO
