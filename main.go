package main

import (
	"backend/db"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
  app := fiber.New()

  app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
  }))

  db.InitDatabase()

  app.Get("/", func (c *fiber.Ctx) error {
    return c.SendString("api is running")
  })

  log.Fatal(app.Listen(":9999"))
}