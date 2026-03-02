USE [maskdump_fixture_multi]
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
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (1, N'i690245@yandex.ru', N'Иван Петров', N'i690245@yandex.ru', N'+7 (506) 405-52-34', 1)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (2, N'(673) 548-0798', N'Emily Carter', N'e526362@yahoo.com', N'(673) 548-0798', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (3, N'lukas.schmidt', N'Lukas Schmidt', N'lba36a4@web.de', N'+45 50 2536 5878', 3)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (4, N'00 32 38 50 01', N'Camille Bernard', N'c404c10@free.fr', N'00 32 38 50 01', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (5, N'ed763e6@telia.se', N'Erik Andersson', N'ed763e6@telia.se', N'+48 7 135 40 65', 1)
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
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (1, N'welcome-playbook', N'Welcome Playbook', N'Escalation contact 1: phone +44 50 4504 1634, email e26860b@news.co.uk. Keep this note in the exported dump.', 1)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (2, N'privacy-checklist', N'Privacy Checklist', N'Escalation contact 2: phone +32 4 58 94 73 72, email p514224@orange.fr. Keep this note in the exported dump.', 2)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (3, N'support-handbook', N'Support Handbook', N'Escalation contact 3: phone +1 (622) 145-4131, email hf31544@gmail.com. Keep this note in the exported dump.', 3)
GO
