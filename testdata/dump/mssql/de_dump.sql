USE [maskdump_fixture_de]
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
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (1, N'lukas.schmidt', N'Lukas Schmidt', N'lukas.schmidt@web.de', N'+49 30 1234 5678', 1)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (2, N'030 123456', N'Anna Muller', N'anna.mueller@gmx.de', N'030 123456', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (3, N'leonie.fischer@mail.de', N'Leonie Fischer', N'leonie.fischer@mail.de', N'+49 (89) 2345 6789', 3)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (4, N'max.weber', N'Max Weber', N'max.weber@t-online.de', N'040 987654', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (5, N'01761234567', N'Sophie Becker', N'sophie.becker@posteo.de', N'01761234567', 1)
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
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (1, N'welcome-playbook', N'Welcome Playbook', N'Escalation contact 1: phone +49 211 4567 8910, email presse@firma.de. Keep this note in the exported dump.', 1)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (2, N'privacy-checklist', N'Privacy Checklist', N'Escalation contact 2: phone 089 998877, email datenschutz@web.de. Keep this note in the exported dump.', 2)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (3, N'support-handbook', N'Support Handbook', N'Escalation contact 3: phone +49 40 7654 3210, email hilfe@gmx.de. Keep this note in the exported dump.', 3)
GO
