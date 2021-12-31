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
  
  var userID = utils.GenID(6)

  var userGenID models.User
  userGenID = models.GetUserByID(userID)
  for userGenID.UserID != "" {
    userID = utils.GenID(6)
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
    Username: body["username"],
    Password: string(hashedPassword),
    People: []string {},
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
    user.Username, 
    user.Email,
    user.People,
  )
    if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // TODO: add an original transaction creation

  return c.JSON(bson.M{
    "_id": user.ID,
    "userID": user.UserID,
    "username": user.Username,
    "email": user.Email,
    "people": user.People,
    "createdAt": user.CreatedAt,
    "updatedAt": user.UpdatedAt,
    "token": tokenString,
  })
}

// @desc   User login
// @route  POST /api/user/login
// @access Public
func loginUser(c *fiber.Ctx) error {
  // getting body
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
    user.Username, 
    user.Email,
    user.People,
  )
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if compareErr == nil {
    return c.JSON(bson.M{
      "_id": user.ID,
      "userID": user.UserID,
      "username": user.Username,
      "email": user.Email,
      "people": user.People,
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

// @desc   Add person to people list
// @route  POST /api/user/people
// @access Private
func addToPeople(c *fiber.Ctx) error {
  // getting body
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  userToAddID := body["userToAddID"]

  usersCollection := db.GetCollection("users")

  var userToAdd models.User
  if err := usersCollection.FindOne(context.Background(), bson.M{
    "userID": userToAddID,
  }).Decode(&userToAdd); err != nil {
    return c.Status(401).JSON(bson.M{
      "message": "Nu există niciun utilizator cu ID-ul introdus.",
    }) 
  }

  // getting the userID
  userIDLocals := fmt.Sprintf("%v", c.Locals("userID"))
  var userID string
  json.Unmarshal([]byte(userIDLocals), &userID)

  var user models.User
  usersCollection.FindOneAndUpdate(context.Background(), bson.M{
    "userID": userID,
  }, bson.M{
    "$push": bson.M{
      "people": userToAddID,
    },
  }).Decode(&user)

  user.People = append(user.People, userToAddID)

  tokenString, err := utils.GenerateToken(
    user.UserID, 
    user.Username, 
    user.Email,
    user.People,
  )
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(bson.M{
    "token": tokenString,
    "people": user.People,
  })
}

// @desc   Get people 
// @route  GET /api/user/people
// @access Private
func getPeople(c *fiber.Ctx) error {
  usersCollection := db.GetCollection("users")
  userIDLocals := fmt.Sprintf("%v", c.Locals("userID"))
  var userID string
  json.Unmarshal([]byte(userIDLocals), &userID)

  // getting the userID
  peopleLocals := fmt.Sprintf("%v", c.Locals("people"))
  var people []string
  json.Unmarshal([]byte(peopleLocals), &people)

  var users []models.User
  // all transactions
  cursor, err := usersCollection.Find(context.Background(), bson.M{
    "userID": bson.M{ 
      "$in": people,
    },
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &users); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(users)
}

// @desc   Change password
// @route  POST /api/user/password
// @access Private
func changePassword(c *fiber.Ctx) error {
  // getting body
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  password := body["password"]
  newPassword := body["newPassword"]

  userIDLocals := fmt.Sprintf("%v", c.Locals("userID"))
  var userID string
  json.Unmarshal([]byte(userIDLocals), &userID)
  
  usersCollection := db.GetCollection("users")

  var user models.User
  usersCollection.FindOne(context.Background(), bson.M{
    "userID": userID,
  }).Decode(&user)

  hashedPassword := user.Password

  compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

  if compareErr == nil {
    hashedPassword, err := bcrypt.GenerateFromPassword(
      []byte(newPassword), 
      10,
    )
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    usersCollection.FindOneAndUpdate(context.Background(), bson.M{
      "userID": userID,
    }, bson.M{
      "$set": bson.M{
        "password": string(hashedPassword),
      },
    }).Decode(&user)
    user.Password = string(hashedPassword)

    return c.JSON(bson.M{
      "_id": user.ID,
      "userID": user.UserID,
      "username": user.Username,
      "hashedPassword": []byte(hashedPassword),
      "email": user.Email,
      "people": user.People,
      "createdAt": user.CreatedAt,
      "updatedAt": user.UpdatedAt,
    })
  } else {
    return c.Status(401).JSON(bson.M{
      "message": "Nu ați introdus parola validă.",
    }) 
  }
}

func getPerson(c *fiber.Ctx) error {
  userID := c.Params("userID")

  usersCollection := db.GetCollection("users")

  var user models.User
  usersCollection.FindOne(context.Background(), bson.M{
    "userID": userID,
  }).Decode(&user)

  return c.JSON(user)
}

