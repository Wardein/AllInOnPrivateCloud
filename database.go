package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialisiere die Datenbank
func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./datenbank.db")
	if err != nil {
		log.Fatal(err)
	}

	// Tabelle erstellen, falls sie nicht existiert
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	/*addAdmin := `
	INSERT INTO users VALUES (
		'1','admin', 'test'
	);`*/
	_, err = db.Exec(createTable)

	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Database initialized successfully.")
}
