package utils

import (
	"math/rand"
	"strconv"
)

func GenID() string {
  var ID string
  for i := 0; i < 6; i++ {
    ID += strconv.Itoa(rand.Intn(9))
  }

  return ID
}