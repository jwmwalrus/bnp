package tests

import (
	"os"
	"testing"

	"github.com/jwmwalrus/bnp/tests/assert"
)

func NewTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "test-bnp-git-*")
	assert.NoError(t, err)
	return dir
}
