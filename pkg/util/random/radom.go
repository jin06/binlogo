package random

import (
	"math/rand"
	"strconv"
	"time"
)

// Uint32 generate a uint32 random
func Uint32() uint32 {
	return uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
}

// String returns a random number in string format
func String() string {
	i := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	return strconv.Itoa(i)
}
