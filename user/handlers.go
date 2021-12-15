package user

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// @desc   Register user
// @route  POST /api/user/register
// @access Public
func registerUser(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  usersCollection := db.GetCollection("users")
  
  var userID = utils.GenID()

  var userGenID models.User
  userGenID = models.GetUserByID(userID)
  for userGenID.UserID != "" {
    userID = utils.GenID()
    userGenID = models.GetUserByID(userID)
  }

  var checkUser models.User
  usersCollection.FindOne(context.Background(), bson.M{
    "email": body["email"],
  }).Decode(&checkUser)
  if checkUser.UserID != "" {
    return c.Status(401).JSON(bson.M{
      "message": "Exista deja un utilizator cu email-ul introdus.",
    })
  }

  hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(body["password"]), 
    10,
  )
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  user := models.User {
    ID: primitive.NewObjectID(),
    UserID: userID,
    Email: body["email"],
    UserName: body["userName"],
    Password: string(hashedPassword),
    Interactions: []string {},
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }

  _, err = usersCollection.InsertOne(context.Background(), user)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // generate token
  tokenString, err := utils.GenerateToken(
    user.UserID,
    user.UserName, 
    user.Email)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // TODO: add an original transaction creation

  return c.JSON(bson.M{
    "_id": user.ID,
    "userID": user.UserID,
    "userName": user.UserName,
    "email": user.Email,
    "interactions": user.Interactions,
    "createdAt": user.CreatedAt,
    "updatedAt": user.UpdatedAt,
    "token": tokenString,
  })
}

// @desc   User login
// @route  POST /api/user/login
// @access Public
func loginUser(c *fiber.Ctx) error {
  // getting doby
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  usersCollection := db.GetCollection("users")

  var user models.User
  if err := usersCollection.FindOne(context.Background(), bson.M{
    "email": body["email"],
  }).Decode(&user); err != nil {
    return c.Status(401).JSON(bson.M{
      "message": "Nu există niciun utilizator cu email-ul introdus.",
    }) 
  }
  hashedPassword := user.Password

  compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body["password"]))

  tokenString, err := utils.GenerateToken(
    user.UserID, 
    user.UserName, 
    user.Email,
  )
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if compareErr == nil {
    return c.JSON(bson.M{
      "_id": user.ID,
      "userID": user.UserID,
      "userName": user.UserName,
      "email": user.Email,
      "interactions": user.Interactions,
      "createdAt": user.CreatedAt,
      "updatedAt": user.UpdatedAt,
      "token": tokenString,
    })
  } else {
    return c.Status(401).JSON(bson.M{
      "message": "Nu ați introdus parola validă.",
    }) 
  }
}