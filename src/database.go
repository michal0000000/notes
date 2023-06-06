package utils

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

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

func FetchYoungestNote(db *sql.DB) string {
	rows, err := db.Query("SELECT MAX(id) FROM notes LIMIT 1")
	if err != nil {
		CheckErr("failed to fetch notes", err)
		return "0"
	}
	defer rows.Close()

	var youngestId *string

	count := 0
	for rows.Next() {
		err := rows.Scan(&youngestId)
		if err != nil {
			CheckErr("failed to scan note", err)
			return "0"
		}
		count += 1
	}

	if count < 2 {
		return "0"
	}

	return *youngestId
}

func UpdateNote(db *sql.DB, noteState *NoteState) error {

	// Prepare query
	stmt, err := db.Prepare("UPDATE notes SET content=? WHERE id=?")
	if err != nil {
		CheckErr("Couldnt update note (1)", err)
		return err
	}

	// Execute transaction
	res, err := stmt.Exec(noteState.Content, int(noteState.Id))

	// Handle errors
	if err != nil {
		CheckErr("Couldnt update note (3)", err)
		return err
	}

	affected, _ := res.RowsAffected()
	fmt.Printf("SAVED! %d\n", affected)
	return nil
}

// Delete note
func DeleteNoteStr(db *sql.DB, noteToDelete string) error {

	// Prepare query
	stmt, err := db.Prepare("DELETE FROM notes WHERE id=?")
	if err != nil {
		CheckErr("Couldnt delete note (1)", err)
		return err
	}

	// Execute query
	noteId, _ := strconv.Atoi(noteToDelete)
	_, err = stmt.Exec(noteId)
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
