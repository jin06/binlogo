package benchmark

import (
	"testing"
)

var s = "skldjfosdf"
var k = 4

func reverseLeftWords(s string, n int) string {
	return s
}

func BenchmarkSubStrRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reverseLeftWords(s, k)
	}
}
