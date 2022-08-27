package ing2

import (
	"crypto/rand"
	"math/big"
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

// GetRandomLetters returns a random string of letters of the given length
func GetRandomLetters(n int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		x, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[x.Int64()]
	}
	return string(b), nil
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
	randomSeed.Store(mrand.New(
		mrand.NewSource(time.Now().UnixNano()),
	))
}
