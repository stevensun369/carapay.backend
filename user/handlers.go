package user

import (
	"github.com/gofiber/fiber/v2"
)

func registerUser(c *fiber.Ctx) error {
  return c.SendString("register user")
}