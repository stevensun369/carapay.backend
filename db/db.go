package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend/env"
)

var Client *mongo.Client

func InitDatabase() {
  var err error
  Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(env.MongoURI))

  if err != nil {
    log.Fatal(err)
  }
}

func GetCollection(collectionName string) (*mongo.Collection) {
  collection := Client.Database("carapay").Collection(collectionName)

  return collection
}