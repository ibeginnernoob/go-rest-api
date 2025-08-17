package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func DBInit() {
	dbCon, err := sql.Open("sqlite3", "api.db?_foreign_keys=on")

	if err != nil {
		panic("db connection could not be established")
	}

	dbCon.SetMaxOpenConns(10)
	dbCon.SetMaxIdleConns(5)

	DB = dbCon

	createTables()
}

func createTables() {
	if DB == nil {
		panic("No db connection")
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,		
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`

	_, err := DB.Exec(createEventsTable)

	if err != nil {
		panic("events table could not be created")
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name STRING TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`

	_, err = DB.Exec(createUsersTable)

	if err != nil {
		panic("users table could not be created")
	}
}
