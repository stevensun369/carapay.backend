package models

import (
	"backend/db"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
  ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  TransactionID string  `json:"transactionID,omitempty" bson:"transactionID,omitempty"`
  From string `json:"from,omitempty" bson:"from,omitempty"`
  To string `json:"to,omitempty" bson:"to,omitempty"`
  Message string `json:"message,omitempty" bson:"message,omitempty"`
  Amount float64 `json:"amount,omitempty" bson:"amount,omitempty"`
  CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
  UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func GetTransactionByID(transactionID string) (Transaction) {
  var transaction Transaction
  transactionsCollection := db.GetCollection("transactionsl")

  transactionsCollection.FindOne(context.Background(), bson.M{
    "transactionID": transactionID,
  }).Decode(&transaction)

  return transaction
}

func CalculateBalance(userID string) (balance float64, err error) {
  var transactionsTo []Transaction
  var transactionsFrom []Transaction
  transactionsCollection := db.GetCollection("transactions")

  // transactionsTo
  cursor, err := transactionsCollection.Find(context.Background(), bson.M{
    "to": userID,
  })
  if err != nil {
    return 0, err
  }
  if err = cursor.All(context.Background(), &transactionsTo); err != nil {
    return 0, err
  }

  // transactionFrom
  cursor, err = transactionsCollection.Find(context.Background(), bson.M{
    "from": userID,
  })
  if err != nil {
    return 0, err
  }
  if err = cursor.All(context.Background(), &transactionsFrom); err != nil {
    return 0, err
  }

  var transactionsToSum float64 = 0
  for _, transaction := range transactionsTo {
    transactionsToSum += transaction.Amount
  }

  var transactionsFromSum float64 = 0
  for _, transaction := range transactionsFrom {
    transactionsFromSum += transaction.Amount
  }

  return transactionsToSum - transactionsFromSum, nil

}