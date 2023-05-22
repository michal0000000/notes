package main

import (
	"github.com/gofiber/fiber/v2" // add Fiber package
	"github.com/gofiber/template/pug"
)

func main() {
	// Initialize Pug template engine in ./views folder
	engine := pug.New("./views", ".html")

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		Views: engine, // set template engine for rendering
	})

	app.Static(
		"/css",         // mount address
		"./static/css", // path to the file folder
	)

	// Create a new endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		// rendering the "index" template with content passing
		return c.Render("index", fiber.Map{
			//...
		})
	})

	// Endpoint for Post method
	app.Post("/", func(c *fiber.Ctx) error {
		// function that stores a new data
		return nil
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
