package utils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateNewNote(db *sql.DB, title string) int64 {
	stmt, err := db.Prepare("INSERT INTO notes(title, content) VALUES (?, ?)")
	if err != nil {
		log.Printf("failed to prepare statement: %v\n", err)
		return 0
	}
	defer stmt.Close()

	res, err := stmt.Exec(title, "")
	if err != nil {
		log.Printf("failed to execute statement: %v\n", err)
		return 0
	}

	noteID, err := res.LastInsertId()
	if err != nil {
		log.Printf("failed to retrieve last inserted ID: %v\n", err)
		return 0
	}

	return noteID
}

func UpdateNote(db *sql.DB, note_id string) error {
	return nil
}

func CreateInitialDatabase(db *sql.DB) {
	// Create notes table
	res, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			content TEXT
		)
	`)
	if err != nil {
		log.Printf("failed to create notes table: %v\n", err)
	}

	// Add more table creation queries if needed
	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("failed to fetch info about created table: %v\n", err)
	}

	// Create default note if new table was created
	if affected > 0 {
		CreateNewNote(db, "Untitled")
	}
}
