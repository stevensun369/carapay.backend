package user

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/api/user")

  // register
  g.Post("/register", registerUser)

  // login
  g.Post("/login", loginUser)
}