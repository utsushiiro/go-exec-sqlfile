package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/format"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

//go:embed db/initdb.d/*
var sqlDir embed.FS

var p *parser.Parser = parser.New()

func TestMain(m *testing.M) {
	InitDB()
	m.Run()
	CloseDB()
}

func ResetAllTables() {
	TruncateAllTables()
	SeedDB()
}

func SeedDB() {
	seedSQL, err := sqlDir.ReadFile("db/initdb.d/02_seed.sql")
	if err != nil {
		log.Fatalf("failed to read seed sql file: %v", err)
	}

	stmts, _, err := p.Parse(string(seedSQL), "", "")
	if err != nil {
		log.Fatalf("failed to parse seed sql: %v", err)
	}

	sqls := []string{}
	for _, stmt := range stmts {
		var buf bytes.Buffer
		stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &buf))
		sqls = append(sqls, buf.String())
	}

	for _, sql := range sqls {
		if _, err := db.Exec(sql); err != nil {
			log.Fatalf("failed to execute sql: err=%v, sql=%s", err, sql)
		}
	}
}

func TruncateAllTables() {
	rows, err := db.Query("show full tables where Table_Type = 'BASE TABLE'")
	if err != nil {
		log.Fatalf("failed to show tables: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, tableType string
		err = rows.Scan(&tableName, &tableType)
		if err != nil {
			log.Fatalf("failed to show tables: %v", err)
		}

		queries := []string{
			"SET FOREIGN_KEY_CHECKS = 0",
			fmt.Sprintf("TRUNCATE `%s`", tableName),
			"SET FOREIGN_KEY_CHECKS = 1",
		}

		for _, q := range queries {
			if _, err := db.Exec(q); err != nil {
				log.Fatalf("failed to truncate table: %#v", err)
			}
		}
	}

	if rows.Err() != nil {
		log.Fatalf("failed to iterate show tables results: %v", rows.Err())
	}
}

func TestResetAllTables(t *testing.T) {
	want := []Post{
		{ID: 1, AuthorID: 1, Content: "Hello, world!"},
		{ID: 2, AuthorID: 1, Content: "(;_;)"},
		{ID: 3, AuthorID: 2, Content: "Hello, world!"},
		{ID: 4, AuthorID: 2, Content: "insert into posts (id, author_id, content) values (4, 2, 'Hello, world!');"},
	}

	got, err := AllPosts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch detected (-want +got):\n%s", diff)
	}

	ResetAllTables()

	got, err = AllPosts(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch detected (-want +got):\n%s", diff)
	}
}