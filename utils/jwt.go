package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTKey = []byte("1234567891011121314")

// claims
type UserClaims struct {
  UserID string `json:"userID"`
  UserName string `json:"userName"`
  Email string `json:"email"`
  jwt.StandardClaims
}

// generations function
func GenerateToken(userID string, userName string, email string) (string, error) {
  // one year has 8760 hours
  expirationTime := time.Now().Add(8760 * time.Hour)

  // the "claims"
  claims := &UserClaims {
    UserID: userID,
    UserName: userName,
    Email: email,
    StandardClaims: jwt.StandardClaims {
      ExpiresAt: expirationTime.Unix(),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err := token.SignedString(JWTKey)

  return tokenString, err
}