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

	// Initialize database
	db, err := sql.Open("sqlite3", "./notes.db")
	utils.CheckErr("Failed to initialize DB", err)

	// Initilize notes array
	allNotes := make(map[int64]utils.Note, 10)

	// Save current note content
	app.Post("/save", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	app.Get("/new_note", func(c *fiber.Ctx) error {
		newNote := utils.Note{
			Id:      utils.CreateNewNote(db, "Untitled"),
			Title:   "Untitled",
			Content: "",
		}
		return c.JSON(newNote)
	})

	app.Use(redirect.New(redirect.Config{
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

	// Endpoint for PUT method
	app.Put("/", func(c *fiber.Ctx) error {
		// function that replaces the existing data
		return nil
	})

	// Endpoint for PATCH method
	app.Patch("/", func(c *fiber.Ctx) error {
		// function that replaces part of the existing data
		return nil
	})

	// Endpoint for DELETE method
	app.Delete("/", func(c *fiber.Ctx) error {
		// function that deletes the data
		return nil
	})

	// Start server on port 3000
	app.Listen(":3000")
}
