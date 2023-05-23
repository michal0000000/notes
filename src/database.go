package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func create_new_note(db *sql.DB, header string) int64 {

	// insert
	stmt, err := db.Prepare("INSERT INTO notes(header, text) values(?,?)")
	checkErr(err)

	res, err := stmt.Exec(header, "")
	checkErr(err)

	noteId, err := res.LastInsertId()
	checkErr(err)

	return noteId
}

func update_note(db *sql.DB, note_id string) error {
	return nil
}
