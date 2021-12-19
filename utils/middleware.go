package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

func UserMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(
    string(authHeader), "Bearer",
  ) {
    token = strings.Fields(string(authHeader))[1]

    // we're parsing the claims
    claims := &UserClaims{}
    tkn, err := jwt.ParseWithClaims(token, claims,
      func (token *jwt.Token) (interface {}, error) {
        return JWTKey, nil
      },
    )

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if !tkn.Valid {
      return c.Status(500).JSON(bson.M{
        "message": "token not valid",
      })
    }

    userIDBytes, _ := json.Marshal(claims.UserID)
    userIDJSON := string(userIDBytes)
    c.Locals("userID", userIDJSON)

    userNameBytes, _ := json.Marshal(claims.UserName)
    userNameJSON := string(userNameBytes)
    c.Locals("userName", userNameJSON)

    emailBytes, _ := json.Marshal(claims.Email)
    emailJSON := string(emailBytes)
    c.Locals("email", emailJSON)
  }

  if (token == "") {
    return c.Status(500).JSON(bson.M{
      "message": "no token",
    })
  }

  return c.Next()
}