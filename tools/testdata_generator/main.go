// Command testdata_generator creates integration fixtures under ../testdata.
//
// Usage:
//
//	cd tools && go run ./testdata_generator
//
// The generator is deterministic. Re-run it after changing fixture content,
// masking rules, or regex presets used by integration tests.
package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	emailRegexPattern = `\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}\b`
	phoneRegexPattern = `(?:(?:\+7|7|8)(?:[\s-]?\(?\d{3}\)?[\s-]?\d{3}[\s-]?\d{2}[\s-]?\d{2}|\d{10})|(?:\+1[\s-]?)?(?:\(?\d{3}\)?[\s-]?\d{3}[\s-]?\d{4})|(?:\+49[\s-]?)?(?:\(?\d{2,4}\)?[\s-]?\d{3,4}[\s-]?\d{3,4})|(?:\+44[\s-]?)?(?:\(?0?\d{2,4}\)?[\s-]?\d{3,4}[\s-]?\d{3,4})|(?:\+33[\s-]?)?(?:\(0\)[\s-]?)?\d(?:[\s-]?\d{2}){4}|(?:\+46[\s-]?)?(?:\(0\)[\s-]?)?\d{1,3}[\s-]?\d{2,3}[\s-]?\d{2}[\s-]?\d{2,3})`
	emailMaskTarget   = "username:2-"
	emailMaskValue    = "hash:6"
	phoneMaskTarget   = "2,3,5,6,8,10"
	phoneMaskValue    = "hash"
)

type group struct {
	ID    int
	Code  string
	Title string
}

type user struct {
	ID      int
	Login   string
	Name    string
	Email   string
	Phone   string
	GroupID int
}

type post struct {
	ID     int
	Code   string
	Title  string
	Detail string
	UserID int
}

type fixture struct {
	Code  string
	Users []user
	Posts []post
}

var (
	emailRegex = regexp.MustCompile(emailRegexPattern)
	phoneRegex = regexp.MustCompile(phoneRegexPattern)
	groups     = []group{
		{ID: 1, Code: "admins", Title: "Administrators"},
		{ID: 2, Code: "editors", Title: "Editorial Team"},
		{ID: 3, Code: "support", Title: "Customer Success"},
	}
	postMeta = []struct {
		ID    int
		Code  string
		Title string
	}{
		{ID: 1, Code: "welcome-playbook", Title: "Welcome Playbook"},
		{ID: 2, Code: "privacy-checklist", Title: "Privacy Checklist"},
		{ID: 3, Code: "support-handbook", Title: "Support Handbook"},
	}
)

func main() {
	root := filepath.Clean(filepath.Join(".."))
	testdataDir := filepath.Join(root, "testdata")

	fixtures := map[string]fixture{
		"ru": {
			Code: "ru",
			Users: []user{
				{ID: 1, Login: "ivan.petrov", Name: "Иван Петров", Email: "ivan.petrov@yandex.ru", Phone: "+7 (916) 555-12-34", GroupID: 1},
				{ID: 2, Login: "8 912 444 55 66", Name: "Анна Смирнова", Email: "anna.smirnova@mail.ru", Phone: "8 912 444 55 66", GroupID: 2},
				{ID: 3, Login: "sergey-volkov", Name: "Сергей Волков", Email: "sergey.volkov@bk.ru", Phone: "7 495 123 45 67", GroupID: 3},
				{ID: 4, Login: "olga.romanova@rambler.ru", Name: "Ольга Романова", Email: "olga.romanova@rambler.ru", Phone: "+7 812 600 77 88", GroupID: 2},
				{ID: 5, Login: "79165550011", Name: "Елена Соколова", Email: "elena.sokolova@list.ru", Phone: "79165550011", GroupID: 1},
			},
			Posts: buildPosts(
				[]string{"+7 (495) 777-11-22", "8 800 555 35 35", "7 812 320 10 10"},
				[]string{"press-office@company.ru", "privacy-team@yandex.ru", "support-center@mail.ru"},
			),
		},
		"us": {
			Code: "us",
			Users: []user{
				{ID: 1, Login: "john.miller", Name: "John Miller", Email: "john.miller@gmail.com", Phone: "+1 (212) 555-0188", GroupID: 1},
				{ID: 2, Login: "(646) 555-0199", Name: "Emily Carter", Email: "emily.carter@yahoo.com", Phone: "(646) 555-0199", GroupID: 2},
				{ID: 3, Login: "mason.hall@outlook.com", Name: "Mason Hall", Email: "mason.hall@outlook.com", Phone: "415-555-0132", GroupID: 3},
				{ID: 4, Login: "olivia.wright", Name: "Olivia Wright", Email: "olivia.wright@proton.me", Phone: "+1 503 555 0114", GroupID: 2},
				{ID: 5, Login: "3125550147", Name: "Noah Davis", Email: "noah.davis@icloud.com", Phone: "3125550147", GroupID: 1},
			},
			Posts: buildPosts(
				[]string{"+1 (202) 555-0141", "415-555-0198", "(646) 555-0102"},
				[]string{"media.desk@newsroom.us", "privacy-team@outlook.com", "help.center@gmail.com"},
			),
		},
		"de": {
			Code: "de",
			Users: []user{
				{ID: 1, Login: "lukas.schmidt", Name: "Lukas Schmidt", Email: "lukas.schmidt@web.de", Phone: "+49 30 1234 5678", GroupID: 1},
				{ID: 2, Login: "030 123456", Name: "Anna Muller", Email: "anna.mueller@gmx.de", Phone: "030 123456", GroupID: 2},
				{ID: 3, Login: "leonie.fischer@mail.de", Name: "Leonie Fischer", Email: "leonie.fischer@mail.de", Phone: "+49 (89) 2345 6789", GroupID: 3},
				{ID: 4, Login: "max.weber", Name: "Max Weber", Email: "max.weber@t-online.de", Phone: "040 987654", GroupID: 2},
				{ID: 5, Login: "01761234567", Name: "Sophie Becker", Email: "sophie.becker@posteo.de", Phone: "01761234567", GroupID: 1},
			},
			Posts: buildPosts(
				[]string{"+49 211 4567 8910", "089 998877", "+49 40 7654 3210"},
				[]string{"presse@firma.de", "datenschutz@web.de", "hilfe@gmx.de"},
			),
		},
		"gb": {
			Code: "gb",
			Users: []user{
				{ID: 1, Login: "oliver.smith", Name: "Oliver Smith", Email: "oliver.smith@btinternet.com", Phone: "+44 20 7946 0958", GroupID: 1},
				{ID: 2, Login: "020 7946 0959", Name: "Amelia Brown", Email: "amelia.brown@outlook.co.uk", Phone: "020 7946 0959", GroupID: 2},
				{ID: 3, Login: "harry.jones@gmail.com", Name: "Harry Jones", Email: "harry.jones@gmail.com", Phone: "+44 161 496 0000", GroupID: 3},
				{ID: 4, Login: "isla.wilson", Name: "Isla Wilson", Email: "isla.wilson@protonmail.com", Phone: "0117 496 0123", GroupID: 2},
				{ID: 5, Login: "07900111222", Name: "George Taylor", Email: "george.taylor@yahoo.co.uk", Phone: "07900111222", GroupID: 1},
			},
			Posts: buildPosts(
				[]string{"+44 113 496 0101", "020 7000 1234", "+44 121 496 0202"},
				[]string{"press.office@news.co.uk", "privacy.unit@outlook.co.uk", "helpdesk@gmail.com"},
			),
		},
		"fr": {
			Code: "fr",
			Users: []user{
				{ID: 1, Login: "luc.martin", Name: "Luc Martin", Email: "luc.martin@orange.fr", Phone: "+33 1 42 68 53 00", GroupID: 1},
				{ID: 2, Login: "01 42 68 53 01", Name: "Camille Bernard", Email: "camille.bernard@free.fr", Phone: "01 42 68 53 01", GroupID: 2},
				{ID: 3, Login: "julie.dubois@sfr.fr", Name: "Julie Dubois", Email: "julie.dubois@sfr.fr", Phone: "+33 (0)4 72 00 00 00", GroupID: 3},
				{ID: 4, Login: "nicolas.moreau", Name: "Nicolas Moreau", Email: "nicolas.moreau@laposte.net", Phone: "06 12 34 56 78", GroupID: 2},
				{ID: 5, Login: "0611223344", Name: "Lea Petit", Email: "lea.petit@proton.me", Phone: "0611223344", GroupID: 1},
			},
			Posts: buildPosts(
				[]string{"+33 1 55 44 33 22", "04 72 10 20 30", "+33 (0)3 88 11 22 33"},
				[]string{"presse@entreprise.fr", "confidentialite@orange.fr", "support-client@free.fr"},
			),
		},
		"se": {
			Code: "se",
			Users: []user{
				{ID: 1, Login: "erik.andersson", Name: "Erik Andersson", Email: "erik.andersson@telia.se", Phone: "+46 8 123 45 67", GroupID: 1},
				{ID: 2, Login: "08-123 45 68", Name: "Anna Johansson", Email: "anna.johansson@outlook.com", Phone: "08-123 45 68", GroupID: 2},
				{ID: 3, Login: "elsa.nilsson@gmail.com", Name: "Elsa Nilsson", Email: "elsa.nilsson@gmail.com", Phone: "+46 (0)31-123 456", GroupID: 3},
				{ID: 4, Login: "oscar.lindberg", Name: "Oscar Lindberg", Email: "oscar.lindberg@bahnhof.se", Phone: "031-765 432", GroupID: 2},
				{ID: 5, Login: "0701234567", Name: "Maja Karlsson", Email: "maja.karlsson@icloud.com", Phone: "0701234567", GroupID: 1},
			},
			Posts: buildPosts(
				[]string{"+46 31 701 23 45", "08-555 12 34", "+46 (0)90-123 456"},
				[]string{"press@bolag.se", "privacy.office@telia.se", "kundservice@bahnhof.se"},
			),
		},
		"multi": {
			Code: "multi",
			Users: []user{
				{ID: 1, Login: "ivan.petrov@yandex.ru", Name: "Иван Петров", Email: "ivan.petrov@yandex.ru", Phone: "+7 (916) 555-12-34", GroupID: 1},
				{ID: 2, Login: "(646) 555-0199", Name: "Emily Carter", Email: "emily.carter@yahoo.com", Phone: "(646) 555-0199", GroupID: 2},
				{ID: 3, Login: "lukas.schmidt", Name: "Lukas Schmidt", Email: "lukas.schmidt@web.de", Phone: "+49 30 1234 5678", GroupID: 3},
				{ID: 4, Login: "01 42 68 53 01", Name: "Camille Bernard", Email: "camille.bernard@free.fr", Phone: "01 42 68 53 01", GroupID: 2},
				{ID: 5, Login: "erik.andersson@telia.se", Name: "Erik Andersson", Email: "erik.andersson@telia.se", Phone: "+46 8 123 45 67", GroupID: 1},
			},
			Posts: buildPosts(
				[]string{"+44 20 7000 1234", "+33 1 55 44 33 22", "+1 (202) 555-0141"},
				[]string{"editorial.office@news.co.uk", "privacy.board@orange.fr", "help.center@gmail.com"},
			),
		},
	}

	sqlRenderers := []struct {
		name   string
		render func(string, fixture) string
	}{
		{name: "mysql", render: renderMySQLDump},
		{name: "oracle", render: renderOracleDump},
		{name: "postgresql", render: renderPostgreSQLDump},
		{name: "mssql", render: renderMSSQLDump},
	}

	for _, renderer := range sqlRenderers {
		for _, code := range []string{"ru", "us", "de", "gb", "fr", "se", "multi"} {
			content := renderer.render(code, fixtures[code])
			inputPath := filepath.Join(testdataDir, "dump", renderer.name, code+"_dump.sql")
			if err := writeTextFile(inputPath, content); err != nil {
				fail(err)
			}

			expectedPath := filepath.Join(testdataDir, "expected", "dump", renderer.name, code+"_dump.sql")
			if err := writeTextFile(expectedPath, maskContent(content)); err != nil {
				fail(err)
			}
		}
	}

	csvContent, err := renderCSV(fixtures["multi"].Users)
	if err != nil {
		fail(err)
	}
	csvPath := filepath.Join(testdataDir, "csv", "tst_users_multi.csv")
	if err := writeTextFile(csvPath, csvContent); err != nil {
		fail(err)
	}

	expectedCSVPath := filepath.Join(testdataDir, "expected", "csv", "tst_users_multi.csv")
	if err := writeTextFile(expectedCSVPath, maskContent(csvContent)); err != nil {
		fail(err)
	}
}

func buildPosts(phones, emails []string) []post {
	posts := make([]post, 0, len(postMeta))
	for i, meta := range postMeta {
		detail := fmt.Sprintf(
			"Escalation contact %d: phone %s, email %s. Keep this note in the exported dump.",
			i+1,
			phones[i],
			emails[i],
		)
		posts = append(posts, post{
			ID:     meta.ID,
			Code:   meta.Code,
			Title:  meta.Title,
			Detail: detail,
			UserID: i + 1,
		})
	}
	return posts
}

func renderMySQLDump(code string, fixture fixture) string {
	var b strings.Builder
	dbName := "maskdump_fixture_" + code

	b.WriteString("-- MySQL dump 10.13  Distrib 8.0.36, for Linux (x86_64)\n")
	b.WriteString("--\n")
	b.WriteString(fmt.Sprintf("-- Host: localhost    Database: %s\n", dbName))
	b.WriteString("-- ------------------------------------------------------\n")
	b.WriteString("/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;\n")
	b.WriteString("/*!40101 SET NAMES utf8mb4 */;\n\n")

	b.WriteString("DROP TABLE IF EXISTS `tst_groups`;\n")
	b.WriteString("CREATE TABLE `tst_groups` (\n")
	b.WriteString("  `id` bigint NOT NULL,\n")
	b.WriteString("  `code` varchar(64) NOT NULL,\n")
	b.WriteString("  `title` varchar(255) NOT NULL,\n")
	b.WriteString("  PRIMARY KEY (`id`)\n")
	b.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n")
	b.WriteString("INSERT INTO `tst_groups` VALUES\n")
	b.WriteString(joinMySQLGroups())
	b.WriteString(";\n\n")

	b.WriteString("DROP TABLE IF EXISTS `tst_users`;\n")
	b.WriteString("CREATE TABLE `tst_users` (\n")
	b.WriteString("  `id` bigint NOT NULL,\n")
	b.WriteString("  `login` varchar(255) NOT NULL,\n")
	b.WriteString("  `name` varchar(255) NOT NULL,\n")
	b.WriteString("  `email` varchar(255) NOT NULL,\n")
	b.WriteString("  `phone` varchar(255) NOT NULL,\n")
	b.WriteString("  `group_id` bigint NOT NULL,\n")
	b.WriteString("  PRIMARY KEY (`id`)\n")
	b.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n")
	b.WriteString("INSERT INTO `tst_users` VALUES\n")
	b.WriteString(joinMySQLUsers(fixture.Users))
	b.WriteString(";\n\n")

	b.WriteString("DROP TABLE IF EXISTS `tst_posts`;\n")
	b.WriteString("CREATE TABLE `tst_posts` (\n")
	b.WriteString("  `id` bigint NOT NULL,\n")
	b.WriteString("  `code` varchar(128) NOT NULL,\n")
	b.WriteString("  `title` varchar(255) NOT NULL,\n")
	b.WriteString("  `detail` text NOT NULL,\n")
	b.WriteString("  `user_id` bigint NOT NULL,\n")
	b.WriteString("  PRIMARY KEY (`id`)\n")
	b.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n")
	b.WriteString("INSERT INTO `tst_posts` VALUES\n")
	b.WriteString(joinMySQLPosts(fixture.Posts))
	b.WriteString(";\n")

	return b.String()
}

func renderPostgreSQLDump(code string, fixture fixture) string {
	var b strings.Builder
	dbName := "maskdump_fixture_" + code

	b.WriteString("--\n")
	b.WriteString("-- PostgreSQL database dump\n")
	b.WriteString("--\n")
	b.WriteString(fmt.Sprintf("-- Database: %s\n", dbName))
	b.WriteString("SET statement_timeout = 0;\n")
	b.WriteString("SET client_encoding = 'UTF8';\n")
	b.WriteString("SET standard_conforming_strings = on;\n\n")

	b.WriteString("CREATE TABLE public.tst_groups (\n")
	b.WriteString("    id bigint PRIMARY KEY,\n")
	b.WriteString("    code varchar(64) NOT NULL,\n")
	b.WriteString("    title varchar(255) NOT NULL\n")
	b.WriteString(");\n")
	b.WriteString("INSERT INTO public.tst_groups (id, code, title) VALUES\n")
	b.WriteString(joinPostgresGroups())
	b.WriteString(";\n\n")

	b.WriteString("CREATE TABLE public.tst_users (\n")
	b.WriteString("    id bigint PRIMARY KEY,\n")
	b.WriteString("    login varchar(255) NOT NULL,\n")
	b.WriteString("    name varchar(255) NOT NULL,\n")
	b.WriteString("    email varchar(255) NOT NULL,\n")
	b.WriteString("    phone varchar(255) NOT NULL,\n")
	b.WriteString("    group_id bigint NOT NULL\n")
	b.WriteString(");\n")
	b.WriteString("INSERT INTO public.tst_users (id, login, name, email, phone, group_id) VALUES\n")
	b.WriteString(joinPostgresUsers(fixture.Users))
	b.WriteString(";\n\n")

	b.WriteString("CREATE TABLE public.tst_posts (\n")
	b.WriteString("    id bigint PRIMARY KEY,\n")
	b.WriteString("    code varchar(128) NOT NULL,\n")
	b.WriteString("    title varchar(255) NOT NULL,\n")
	b.WriteString("    detail text NOT NULL,\n")
	b.WriteString("    user_id bigint NOT NULL\n")
	b.WriteString(");\n")
	b.WriteString("INSERT INTO public.tst_posts (id, code, title, detail, user_id) VALUES\n")
	b.WriteString(joinPostgresPosts(fixture.Posts))
	b.WriteString(";\n")

	return b.String()
}

func renderMSSQLDump(code string, fixture fixture) string {
	var b strings.Builder
	dbName := "maskdump_fixture_" + code

	b.WriteString("USE [" + dbName + "]\n")
	b.WriteString("GO\n")
	b.WriteString("SET ANSI_NULLS ON\n")
	b.WriteString("GO\n")
	b.WriteString("SET QUOTED_IDENTIFIER ON\n")
	b.WriteString("GO\n\n")

	b.WriteString("CREATE TABLE [dbo].[tst_groups](\n")
	b.WriteString("    [id] bigint NOT NULL,\n")
	b.WriteString("    [code] nvarchar(64) NOT NULL,\n")
	b.WriteString("    [title] nvarchar(255) NOT NULL,\n")
	b.WriteString("    CONSTRAINT [PK_tst_groups] PRIMARY KEY CLUSTERED ([id] ASC)\n")
	b.WriteString(")\n")
	b.WriteString("GO\n")
	for _, item := range groups {
		b.WriteString(fmt.Sprintf(
			"INSERT INTO [dbo].[tst_groups] ([id], [code], [title]) VALUES (%d, N'%s', N'%s')\n",
			item.ID,
			escapeSQL(item.Code),
			escapeSQL(item.Title),
		))
	}
	b.WriteString("GO\n\n")

	b.WriteString("CREATE TABLE [dbo].[tst_users](\n")
	b.WriteString("    [id] bigint NOT NULL,\n")
	b.WriteString("    [login] nvarchar(255) NOT NULL,\n")
	b.WriteString("    [name] nvarchar(255) NOT NULL,\n")
	b.WriteString("    [email] nvarchar(255) NOT NULL,\n")
	b.WriteString("    [phone] nvarchar(255) NOT NULL,\n")
	b.WriteString("    [group_id] bigint NOT NULL,\n")
	b.WriteString("    CONSTRAINT [PK_tst_users] PRIMARY KEY CLUSTERED ([id] ASC)\n")
	b.WriteString(")\n")
	b.WriteString("GO\n")
	for _, item := range fixture.Users {
		b.WriteString(fmt.Sprintf(
			"INSERT INTO [dbo].[tst_users] ([id], [login], [name], [email], [phone], [group_id]) VALUES (%d, N'%s', N'%s', N'%s', N'%s', %d)\n",
			item.ID,
			escapeSQL(item.Login),
			escapeSQL(item.Name),
			escapeSQL(item.Email),
			escapeSQL(item.Phone),
			item.GroupID,
		))
	}
	b.WriteString("GO\n\n")

	b.WriteString("CREATE TABLE [dbo].[tst_posts](\n")
	b.WriteString("    [id] bigint NOT NULL,\n")
	b.WriteString("    [code] nvarchar(128) NOT NULL,\n")
	b.WriteString("    [title] nvarchar(255) NOT NULL,\n")
	b.WriteString("    [detail] nvarchar(max) NOT NULL,\n")
	b.WriteString("    [user_id] bigint NOT NULL,\n")
	b.WriteString("    CONSTRAINT [PK_tst_posts] PRIMARY KEY CLUSTERED ([id] ASC)\n")
	b.WriteString(")\n")
	b.WriteString("GO\n")
	for _, item := range fixture.Posts {
		b.WriteString(fmt.Sprintf(
			"INSERT INTO [dbo].[tst_posts] ([id], [code], [title], [detail], [user_id]) VALUES (%d, N'%s', N'%s', N'%s', %d)\n",
			item.ID,
			escapeSQL(item.Code),
			escapeSQL(item.Title),
			escapeSQL(item.Detail),
			item.UserID,
		))
	}
	b.WriteString("GO\n")

	return b.String()
}

func renderOracleDump(code string, fixture fixture) string {
	var b strings.Builder
	schemaName := "MASKDUMP_" + strings.ToUpper(code)

	b.WriteString("-- Oracle Database dump\n")
	b.WriteString(fmt.Sprintf("-- Schema: %s\n", schemaName))
	b.WriteString("SET DEFINE OFF;\n\n")
	b.WriteString("BEGIN EXECUTE IMMEDIATE 'DROP TABLE tst_groups'; EXCEPTION WHEN OTHERS THEN NULL; END;\n/\n")
	b.WriteString("BEGIN EXECUTE IMMEDIATE 'DROP TABLE tst_users'; EXCEPTION WHEN OTHERS THEN NULL; END;\n/\n")
	b.WriteString("BEGIN EXECUTE IMMEDIATE 'DROP TABLE tst_posts'; EXCEPTION WHEN OTHERS THEN NULL; END;\n/\n\n")

	b.WriteString("CREATE TABLE tst_groups (\n")
	b.WriteString("    id NUMBER(19) PRIMARY KEY,\n")
	b.WriteString("    code VARCHAR2(64 CHAR) NOT NULL,\n")
	b.WriteString("    title VARCHAR2(255 CHAR) NOT NULL\n")
	b.WriteString(");\n")
	for _, item := range groups {
		b.WriteString(fmt.Sprintf(
			"INSERT INTO tst_groups (id, code, title) VALUES (%d, '%s', '%s');\n",
			item.ID,
			escapeSQL(item.Code),
			escapeSQL(item.Title),
		))
	}
	b.WriteString("\n")

	b.WriteString("CREATE TABLE tst_users (\n")
	b.WriteString("    id NUMBER(19) PRIMARY KEY,\n")
	b.WriteString("    login VARCHAR2(255 CHAR) NOT NULL,\n")
	b.WriteString("    name VARCHAR2(255 CHAR) NOT NULL,\n")
	b.WriteString("    email VARCHAR2(255 CHAR) NOT NULL,\n")
	b.WriteString("    phone VARCHAR2(255 CHAR) NOT NULL,\n")
	b.WriteString("    group_id NUMBER(19) NOT NULL\n")
	b.WriteString(");\n")
	for _, item := range fixture.Users {
		b.WriteString(fmt.Sprintf(
			"INSERT INTO tst_users (id, login, name, email, phone, group_id) VALUES (%d, '%s', '%s', '%s', '%s', %d);\n",
			item.ID,
			escapeSQL(item.Login),
			escapeSQL(item.Name),
			escapeSQL(item.Email),
			escapeSQL(item.Phone),
			item.GroupID,
		))
	}
	b.WriteString("\n")

	b.WriteString("CREATE TABLE tst_posts (\n")
	b.WriteString("    id NUMBER(19) PRIMARY KEY,\n")
	b.WriteString("    code VARCHAR2(128 CHAR) NOT NULL,\n")
	b.WriteString("    title VARCHAR2(255 CHAR) NOT NULL,\n")
	b.WriteString("    detail CLOB NOT NULL,\n")
	b.WriteString("    user_id NUMBER(19) NOT NULL\n")
	b.WriteString(");\n")
	for _, item := range fixture.Posts {
		b.WriteString(fmt.Sprintf(
			"INSERT INTO tst_posts (id, code, title, detail, user_id) VALUES (%d, '%s', '%s', '%s', %d);\n",
			item.ID,
			escapeSQL(item.Code),
			escapeSQL(item.Title),
			escapeSQL(item.Detail),
			item.UserID,
		))
	}
	b.WriteString("COMMIT;\n")

	return b.String()
}

func renderCSV(users []user) (string, error) {
	var b strings.Builder
	writer := csv.NewWriter(&b)

	if err := writer.Write([]string{"id", "login", "name", "email", "phone", "group_id"}); err != nil {
		return "", err
	}

	for _, item := range users {
		record := []string{
			strconv.Itoa(item.ID),
			item.Login,
			item.Name,
			item.Email,
			item.Phone,
			strconv.Itoa(item.GroupID),
		}
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return b.String(), nil
}

func joinMySQLGroups() string {
	parts := make([]string, 0, len(groups))
	for _, item := range groups {
		parts = append(parts, fmt.Sprintf("(%d,'%s','%s')", item.ID, escapeSQL(item.Code), escapeSQL(item.Title)))
	}
	return strings.Join(parts, ",\n")
}

func joinMySQLUsers(users []user) string {
	parts := make([]string, 0, len(users))
	for _, item := range users {
		parts = append(parts, fmt.Sprintf(
			"(%d,'%s','%s','%s','%s',%d)",
			item.ID,
			escapeSQL(item.Login),
			escapeSQL(item.Name),
			escapeSQL(item.Email),
			escapeSQL(item.Phone),
			item.GroupID,
		))
	}
	return strings.Join(parts, ",\n")
}

func joinMySQLPosts(posts []post) string {
	parts := make([]string, 0, len(posts))
	for _, item := range posts {
		parts = append(parts, fmt.Sprintf(
			"(%d,'%s','%s','%s',%d)",
			item.ID,
			escapeSQL(item.Code),
			escapeSQL(item.Title),
			escapeSQL(item.Detail),
			item.UserID,
		))
	}
	return strings.Join(parts, ",\n")
}

func joinPostgresGroups() string {
	parts := make([]string, 0, len(groups))
	for _, item := range groups {
		parts = append(parts, fmt.Sprintf("(%d, '%s', '%s')", item.ID, escapeSQL(item.Code), escapeSQL(item.Title)))
	}
	return strings.Join(parts, ",\n")
}

func joinPostgresUsers(users []user) string {
	parts := make([]string, 0, len(users))
	for _, item := range users {
		parts = append(parts, fmt.Sprintf(
			"(%d, '%s', '%s', '%s', '%s', %d)",
			item.ID,
			escapeSQL(item.Login),
			escapeSQL(item.Name),
			escapeSQL(item.Email),
			escapeSQL(item.Phone),
			item.GroupID,
		))
	}
	return strings.Join(parts, ",\n")
}

func joinPostgresPosts(posts []post) string {
	parts := make([]string, 0, len(posts))
	for _, item := range posts {
		parts = append(parts, fmt.Sprintf(
			"(%d, '%s', '%s', '%s', %d)",
			item.ID,
			escapeSQL(item.Code),
			escapeSQL(item.Title),
			escapeSQL(item.Detail),
			item.UserID,
		))
	}
	return strings.Join(parts, ",\n")
}

func maskContent(content string) string {
	parts := strings.SplitAfter(content, "\n")
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		part = emailRegex.ReplaceAllStringFunc(part, maskEmail)
		part = phoneRegex.ReplaceAllStringFunc(part, maskPhone)
		if strings.TrimSpace(part) != "" {
			filtered = append(filtered, part)
		}
	}
	return strings.Join(filtered, "")
}

func maskEmail(email string) string {
	segments := strings.Split(email, "@")
	if len(segments) != 2 {
		return email
	}

	localPart := segments[0]
	domainPart := segments[1]
	positions := parseTargetPositions(strings.TrimPrefix(emailMaskTarget, "username:"), len(localPart))

	return applyMasking(localPart, positions, emailMaskValue, true) + "@" + domainPart
}

func maskPhone(phone string) string {
	digitsOnly := extractDigits(phone)
	positions := parseTargetPositions(phoneMaskTarget, len(digitsOnly))
	maskedDigits := applyMasking(digitsOnly, positions, phoneMaskValue, false)

	var out strings.Builder
	digitIndex := 0
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			out.WriteByte(maskedDigits[digitIndex])
			digitIndex++
			continue
		}
		out.WriteRune(r)
	}

	return out.String()
}

func parseTargetPositions(target string, length int) []int {
	var positions []int

	if strings.Contains(target, "-") {
		parts := strings.Split(target, "-")
		start := 1
		end := length

		if parts[0] != "" {
			start, _ = strconv.Atoi(parts[0])
		}
		if parts[1] != "" {
			end, _ = strconv.Atoi(parts[1])
		}

		for i := start; i <= end && i <= length; i++ {
			positions = append(positions, i-1)
		}
		return positions
	}

	for _, item := range strings.Split(target, ",") {
		pos, _ := strconv.Atoi(item)
		if pos > 0 && pos <= length {
			positions = append(positions, pos-1)
		}
	}

	return positions
}

func applyMasking(value string, positions []int, maskValue string, isEmail bool) string {
	runes := []rune(value)
	maskRunes := []rune{}
	hash := ""

	if maskValue == "*" {
		maskRunes = make([]rune, len(positions))
		for i := range maskRunes {
			maskRunes[i] = '*'
		}
	} else if strings.HasPrefix(maskValue, "hash") {
		hashLen := 16
		if strings.HasPrefix(maskValue, "hash:") {
			parts := strings.Split(maskValue, ":")
			hashLen, _ = strconv.Atoi(parts[1])

			sum := md5.Sum([]byte(value))
			hash = hex.EncodeToString(sum[:])[:hashLen]
		} else if isEmail {
			sum := md5.Sum([]byte(value))
			hash = hex.EncodeToString(sum[:])[:len(runes)]
		} else {
			sum := sha256.Sum256([]byte(value))
			hash = extractDigits(hex.EncodeToString(sum[:]))
		}

		maskRunes = []rune(hash)
	}

	if isEmail && strings.HasPrefix(maskValue, "hash:") && isContinuousSequence(positions) {
		return replacePositions(value, positions, hash)
	}

	for i, pos := range positions {
		if pos >= 0 && pos < len(runes) && i < len(maskRunes) {
			runes[pos] = maskRunes[i]
		}
	}

	return string(runes)
}

func isContinuousSequence(positions []int) bool {
	if len(positions) <= 1 {
		return true
	}
	for i := 1; i < len(positions); i++ {
		if positions[i] != positions[i-1]+1 {
			return false
		}
	}
	return true
}

func replacePositions(value string, positions []int, hash string) string {
	if len(positions) == 0 {
		return value
	}

	runes := []rune(value)
	var result []rune

	prev := 0
	for _, pos := range positions {
		if pos < 0 || pos >= len(runes) {
			continue
		}
		result = append(result, runes[prev:pos]...)
		prev = pos + 1
	}
	result = append(result, runes[prev:]...)

	insertPos := positions[0]
	if insertPos < 0 {
		insertPos = 0
	}
	if insertPos > len(result) {
		insertPos = len(result)
	}

	final := make([]rune, 0, len(result)+len(hash))
	final = append(final, result[:insertPos]...)
	final = append(final, []rune(hash)...)
	final = append(final, result[insertPos:]...)

	return string(final)
}

func extractDigits(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func escapeSQL(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func writeTextFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func fail(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "testdata_generator: %v\n", err)
	os.Exit(1)
}
