package user

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/api/user")

  g.Post("/register", registerUser)
}