package utils

import (
	"database/sql"
)

type Note struct {
	Id      int64
	Title   string
	Content string
}

func FetchNotes(db *sql.DB, allNotes map[int64]Note, amount int) int {
	// Query first 10 notes
	// Select only Id,Title to optimize performance
	rows, err := db.Query("SELECT Id,Title FROM notes LIMIT ?", amount)
	if err != nil {
		CheckErr("failed to fetch notes", err)
		return 0
	}
	defer rows.Close()

	// Append notes
	var count = 0
	for rows.Next() {
		newNote := Note{}
		err := rows.Scan(&newNote.Id, &newNote.Title)
		if err != nil {
			CheckErr("failed to scan note", err)
			return count
		}

		// Append note to allNotes if it doesnt exist yet
		if _, ok := allNotes[newNote.Id]; !ok {
			note := Note{}
			note.Id = newNote.Id
			note.Title = newNote.Title
			note.Content = "" // Make content empty for now
			allNotes[newNote.Id] = note
		}

		count += 1
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		CheckErr("error occurred during iteration", err)
		return count
	}

	return count
}

func FetchSingleNote(db *sql.DB, allNotes map[int64]Note, noteId int64) (string, error) {
	// Query first 10 notes
	rows, err := db.Query("SELECT content FROM notes WHERE id=?", noteId)
	if err != nil {
		CheckErr("failed to fetch notes", err)
		return "", err
	}
	defer rows.Close()

	// Append notes
	var content = ""
	for rows.Next() {
		err := rows.Scan(&content)
		if err != nil {
			CheckErr("failed to scan note", err)
			return "", err
		}

		// Append content of note to allNotes
		note := allNotes[noteId]
		note.Content = content // Make content empty for now
		allNotes[noteId] = note
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		CheckErr("error occurred during iteration", err)
		return "", err
	}

	return content, err

}
