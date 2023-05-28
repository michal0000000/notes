package utils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateNewNote(db *sql.DB, title string) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO notes(title, content) VALUES (?, ?)")
	if err != nil {
		log.Printf("failed to prepare statement: %v\n", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(title, "")
	if err != nil {
		log.Printf("failed to execute statement: %v\n", err)
		return 0, err
	}

	noteID, err := res.LastInsertId()
	if err != nil {
		log.Printf("failed to retrieve last inserted ID: %v\n", err)
		return 0, err
	}

	return noteID, err
}

func UpdateNote(db *sql.DB, noteState *NoteState) error {

	// Prepare query
	updateSQL, err := db.Prepare("UPDATE notes SET content=? where id=?")
	if err != nil {
		CheckErr("Couldnt update note (1)", err)
		return err
	}
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		CheckErr("Couldnt update note (2)", err)
		return err
	}

	// Execute transaction
	_, err = tx.Stmt(updateSQL).Exec(noteState.Content, noteState.Id)

	// Handle errors
	if err != nil {
		CheckErr("Couldnt update note (3)", err)
		log.Println("Doing rollback")
		tx.Rollback()
		return err
	} else {
		tx.Commit()
	}

	return nil
}

// Delete note
func DeleteNote(db *sql.DB, noteState *NoteState) error {

	// Prepare query
	stmt, err := db.Prepare("DELETE FROM notes WHERE id=?")
	if err != nil {
		CheckErr("Couldnt delete note (1)", err)
		return err
	}

	// Execute query
	_, err = stmt.Exec(noteState.Id)
	if err != nil {
		CheckErr("Couldnt delete note (2)", err)
		return err
	}

	return err
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
