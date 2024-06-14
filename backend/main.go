package main

import (
	"lru-cache/server"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Post("/", server.PutCache)
	app.Get("/snapshot", server.GetSnapshot)
	app.Delete("/clear", server.ClearCache)
	app.Get("/:key", server.GetCache)
	app.Delete("/:key", server.DeleteCache)
	app.Listen(":3000")
}
