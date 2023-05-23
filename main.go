package main

import (
	"encoding/json"

	utils "notes/src"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug"
)

func main() {
	// Initialize Pug template engine in ./views folder
	engine := pug.New("./views", ".html")

	// Create a new Fiber instancse
	app := fiber.New(fiber.Config{
		Views: engine, // set template engine for rendering
	})

	// mount // path to folder
	app.Static("/css", "./static/css")
	app.Static("/js", "./static/js")

	// Create a new endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		// rendering the "index" template with content passing
		return c.Render("index", fiber.Map{
			//...
		})
	})

	app.Post("/save", func(c *fiber.Ctx) error {
		// function that stores a new data
		return nil
	})

	app.Post("/new_note", func(c *fiber.Ctx) error {
		newNote := utils.Note{
			Id:      2342,
			Header:  "Quotes",
			Content: "<h2>I am not the body, I am not even the mind</h2>",
		}
		n, err := json.Marshal(newNote)
		if err != nil {
			panic(err)
		}
		return c.JSON(string(n))
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
