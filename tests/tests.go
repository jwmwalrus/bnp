package tests

import (
	"os"
	"testing"
)

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func NewTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "test-bnp-git-*")
	AssertNoError(t, err)
	return dir
}
