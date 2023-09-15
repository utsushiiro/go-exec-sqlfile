package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	err := InitDB()
	if err != nil {
		panic(err)
	}
	defer CloseDB()
}

func InitDB() error {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	config := mysql.Config{
		DBName:    "test",
		User:      "root",
		Passwd:    "password",
		Addr:      "localhost:3306",
		Net:       "tcp",
		Collation: "utf8mb4_bin",
		ParseTime: true,
		Loc:       jst,
	}

	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}
}
