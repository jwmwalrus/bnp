package ing2

import (
	mrand "math/rand"
	"sync/atomic"
	"time"
)

var randomSeed atomic.Pointer[mrand.Rand]

// GetRandomString returns a random string of the given length
func GetRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[randomSeed.Load().Intn(len(charset))]
	}
	return string(b)
}

// TruncateText returns the truncated version of the given string,
// ellipsis added
func TruncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:max] + "..."
}

func init() {
	randomSeed = rand.New(
		rand.NewSource(time.Now().UnixNano()))
}
