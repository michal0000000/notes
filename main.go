package main

import (
	"database/sql"
	"log"
	utils "notes/src"
	"os"
	"strconv"

	"github.com/gofiber/template/html"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/redirect"
)

func init() {
	// Set log output
	log.SetOutput(os.Stdout)

	// Initialize database
	db, err := sql.Open("sqlite3", "./notes.db")
	defer db.Close()
	utils.CheckErr("Failed to initialize DB in init", err)
	utils.CreateInitialDatabase(db)
}

func main() {

	// Initialize Pug template engine in ./views folder
	engine := html.New("./views", ".html")

	// Create a new Fiber instancse
	app := fiber.New(fiber.Config{
		Views: engine, // set template engine for rendering
	})

	// mount // path to folder
	app.Static("/css", "./static/css")
	app.Static("/js", "./static/js")
	app.Static("/img", "./static/img")

	// Initialize database
	db, err := sql.Open("sqlite3", "./notes.db")
	defer db.Close()
	utils.CheckErr("Failed to initialize DB", err)

	// Initilize notes array
	allNotes := make(map[int64]utils.Note, 10)

	// Save current note content
	app.Put("/save", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		// Parse current note state from request data
		currentNoteState := new(utils.NoteState)
		err := c.BodyParser(currentNoteState)
		utils.CheckErr("could not parse note state", err)
		if err != nil {
			return c.JSON(fiber.Map{
				"message": "failed to update note",
			})
		}

		utils.UpdateNote(db, currentNoteState)

		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	// Creat new note
	app.Get("/new_note", func(c *fiber.Ctx) error {

		newNoteId, err := utils.CreateNewNote(db, "Untitled")
		if err != nil {
			return c.JSON(fiber.Map{
				"message": "failed to create new note",
			})
		}

		newNote := utils.Note{
			Id:      newNoteId,
			Title:   "Untitled",
			Content: "",
		}
		return c.JSON(newNote)
	})

	// Redirect root to some specific note
	app.Use(redirect.New(redirect.Config{
		// TODO: doesnt work
		Rules: map[string]string{
			"/": "/:noteId",
		},
		StatusCode: 301,
	}))

	app.Get("/:noteId", func(c *fiber.Ctx) error {

		// Fetch first 10 Notes
		utils.FetchNotes(db, allNotes, 10)

		// Fetch content of selected noteId
		noteId, err := strconv.ParseInt(c.Params("noteId"), 10, 64)
		if err != nil {
			utils.CheckErr("invalid nodeId", err)
			// TODO: redirect to '/'
		}

		selectedNote, err := utils.FetchSingleNote(db, allNotes, noteId)
		if err != nil {
			utils.CheckErr("couldnt retrieve note content", err)
			// TODO: redirect to '/'
		}

		// rendering the "index" template with content passing
		return c.Render("index", fiber.Map{
			"AllNotes":     allNotes,
			"SelectedNote": selectedNote,
		})
	})

	// Endpoint for DELETE method
	app.Delete("/delete", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		// Parse current note state from request data
		noteToDelete := new(utils.NoteState)
		err := c.BodyParser(noteToDelete)
		utils.CheckErr("could not parse note to delete", err)
		if err != nil {
			return c.JSON(fiber.Map{
				"message": "failed to delete note",
			})
		}

		utils.DeleteNote(db, noteToDelete)

		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	// Start server on port 3000
	app.Listen("192.168.1.4:3000")
}
