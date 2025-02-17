package storage

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "api.db")

	if err != nil {
		panic("could not connect to databse")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`
	_, err := DB.Exec(createUsersTable)

	if err != nil {
		panic(err)
	}

	createRoomsTable := `
	CREATE TABLE IF NOT EXISTS rooms (
    	id TEXT PRIMARY KEY,
    	name TEXT NOT NULL,
		description TEXT NULL,
    	host_id INTEGER NOT NULL,
    	public BOOLEAN NOT NULL DEFAULT true,
    	created_at DATETIME NOT NULL,
    	FOREIGN KEY (host_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createRoomsTable)

	if err != nil {
		panic(err)
	}
}
