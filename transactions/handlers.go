package transactions

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func getTransactions(c *fiber.Ctx) error {
  // getting the userID
  userIDLocals := fmt.Sprintf("%v", c.Locals("userID"))
  var userID string
  json.Unmarshal([]byte(userIDLocals), &userID)
  fmt.Println("got the userID")

  // var transactionsTo []models.Transaction
  // var transactionsFrom []models.Transaction
  var transactions []models.Transaction
  transactionsCollection := db.GetCollection("transactions")
  fmt.Println("got the transactions collection")

  // // transactionsTo
  // cursor, err := transactionsCollection.Find(context.Background(), bson.M{
  //   "to": userID,
  // })
  // if err != nil {
  //   return c.Status(500).SendString(fmt.Sprintf("%v", err))
  // }
  // if err = cursor.All(context.Background(), &transactionsTo); err != nil {
  //   return c.Status(500).SendString(fmt.Sprintf("%v", err))
  // }

  // // transactionFrom
  // cursor, err = transactionsCollection.Find(context.Background(), bson.M{
  //   "from": userID,
  // })
  // if err != nil {
  //   return c.Status(500).SendString(fmt.Sprintf("%v", err))
  // }
  // if err = cursor.All(context.Background(), &transactionsFrom); err != nil {
  //   return c.Status(500).SendString(fmt.Sprintf("%v", err))
  // }

  // all transactions
  cursor, err := transactionsCollection.Find(context.Background(), bson.M{
    "$or": []bson.M {
      bson.M{"to": userID},
      bson.M{"from": userID},
    },
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &transactions); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // transactions := append(transactionsTo, transactionsFrom...)

  return c.JSON(transactions)
}

func getTransaction(c *fiber.Ctx) error {
  transactionID := c.Params("transactionID")

  transactionsCollection := db.GetCollection("transactions")

  var transaction models.Transaction
  transactionsCollection.FindOne(context.Background(), bson.M{
    "transactionID": transactionID,
  }).Decode(&transaction)

  return c.JSON(transaction)
}

func createTransaction(c *fiber.Ctx) error {
  // get body
  var body map[string]string
  json.Unmarshal([]byte(c.Body()), &body)
  
  // get body values
  to := body["to"]
  message := body["message"]
  amount := body["amount"]
  password := body["password"]
  
  // get id
  var transactionID = utils.GenID(12)
  var transactionGenID models.Transaction
  transactionGenID = models.GetTransactionByID(transactionID)
  for transactionGenID.TransactionID != "" {
    transactionID = utils.GenID(12)
    transactionGenID = models.GetTransactionByID(transactionID)
  }
  
  // getting the userID
  userIDLocals := fmt.Sprintf("%v", c.Locals("userID"))
  var userID string
  json.Unmarshal([]byte(userIDLocals), &userID)

  // checking if the amount is smaller than user's balance
  balance, err := models.CalculateBalance(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  amountFloat, _ := strconv.ParseFloat(amount, 8)
  if amountFloat > balance {
    return c.Status(401).JSON(bson.M{
      "message": "Fonduri indisponibile.",
    })
  }

  // get user
  var user models.User
  usersCollection := db.GetCollection("users")
  usersCollection.FindOne(context.Background(), bson.M{
    "userID": userID,
  }).Decode(&user)

  // transactions collection
  transactionsCollection := db.GetCollection("transactions")

  // check password
  hashedPassword := user.Password

  compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

  if compareErr == nil {
    transaction := models.Transaction {
      ID: primitive.NewObjectID(),
      TransactionID: transactionID,
      From: userID,
      To: to,
      Message: message,
      Amount: amountFloat,
      CreatedAt: time.Now(),
      UpdatedAt: time.Now(),
    }

    _, err := transactionsCollection.InsertOne(context.Background(), transaction)
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    return c.JSON(transaction)
  } else {
    return c.Status(401).JSON(bson.M{
      "message": "Nu ați introdus parola validă.",
    })
  }
}

func getBalance(c *fiber.Ctx) error {
  // getting the userID
  userIDLocals := fmt.Sprintf("%v", c.Locals("userID"))
  var userID string
  json.Unmarshal([]byte(userIDLocals), &userID)

  balance, err := models.CalculateBalance(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(bson.M{
    "balance": balance,
  })
}

