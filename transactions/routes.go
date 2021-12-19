package transactions

import (
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/api/transactions")

  // get all transactoins
  g.Get("/", utils.UserMiddleware, getTransactions)

  // get transaction
  g.Get("/transaction/:transactionID", utils.UserMiddleware, getTransaction)

  // create transaction
  g.Post("/", utils.UserMiddleware, createTransaction)

  // get balance
  g.Get("/balance", utils.UserMiddleware, getBalance)
}