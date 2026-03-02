USE [maskdump_fixture_us]
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
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (1, N'john.miller', N'John Miller', N'john.miller@gmail.com', N'+1 (212) 555-0188', 1)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (2, N'(646) 555-0199', N'Emily Carter', N'emily.carter@yahoo.com', N'(646) 555-0199', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (3, N'mason.hall@outlook.com', N'Mason Hall', N'mason.hall@outlook.com', N'415-555-0132', 3)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (4, N'olivia.wright', N'Olivia Wright', N'olivia.wright@proton.me', N'+1 503 555 0114', 2)
INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (5, N'3125550147', N'Noah Davis', N'noah.davis@icloud.com', N'3125550147', 1)
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
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (1, N'welcome-playbook', N'Welcome Playbook', N'Escalation contact 1: phone +1 (202) 555-0141, email media.desk@newsroom.us. Keep this note in the exported dump.', 1)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (2, N'privacy-checklist', N'Privacy Checklist', N'Escalation contact 2: phone 415-555-0198, email privacy-team@outlook.com. Keep this note in the exported dump.', 2)
INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (3, N'support-handbook', N'Support Handbook', N'Escalation contact 3: phone (646) 555-0102, email help.center@gmail.com. Keep this note in the exported dump.', 3)
GO
