package utility

import (
	"math/rand"
	"time"
)

func GetRandomNumbersInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
