package utils

import (
	"math/rand"
	"strconv"
)

func GenID(charnum int) string {
  var ID string
  for i := 0; i < charnum; i++ {
    ID += strconv.Itoa(rand.Intn(9))
  }

  return ID
}