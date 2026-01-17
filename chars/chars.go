package chars

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
	"sync/atomic"
	"time"
)

var randomSeed atomic.Pointer[mrand.Rand]

// GetRandomString returns a random string of the given length.
func GetRandomString(n int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	return randomFromCharset(charset, n)
}

// GetRandomLetters returns a random string of letters of the given length.
func GetRandomLetters(n int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	return randomFromCharset(charset, n)
}

func randomFromCharset(charset string, n int) (string, error) {
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

// Truncate returns the truncated version of the given string,
// ellipsis added.
func Truncate(s string, n int, noEllipsis ...bool) string {
	if len(s) <= n {
		return s
	}
	if len(noEllipsis) > 0 && noEllipsis[0] {
		return s[:n]
	}
	return s[:max(0, n-3)] + "..."
}

func init() {
	randomSeed.Store(mrand.New(
		mrand.NewSource(time.Now().UnixNano()),
	))
}
