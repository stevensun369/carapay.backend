package user

import (
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/api/user")

  // register
  g.Post("/register", registerUser)

  // login
  g.Post("/login", loginUser)

  // change password
  g.Post("/password", utils.UserMiddleware, changePassword)

  // add to pelple
  g.Post("/people", utils.UserMiddleware, addToPeople)

  // get people
  g.Get("/people", utils.UserMiddleware, getPeople)

}