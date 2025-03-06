package random

import (
	"math/rand"
	"time"
)

// Uint32 generate a uint32 random
func Uint32() uint32 {
	return uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
}
