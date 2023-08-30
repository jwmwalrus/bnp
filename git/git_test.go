package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jwmwalrus/bnp/tests"
)

func TestAddToStaging(t *testing.T) {
	dir := tests.NewTempDir(t)
	defer os.RemoveAll(dir)

	err := os.Chdir(dir)
	tests.AssertNoError(t, err)

	g, err := NewInterface(dir)
	tests.AssertNoError(t, err)

	err = g.Init()
	tests.AssertNoError(t, err)

	file := "test-file.txt"

	err = os.WriteFile(filepath.Join(dir, file), []byte("test"), 0644)
	tests.AssertNoError(t, err)

	err = g.AddToStaging([]string{file})
	tests.AssertNoError(t, err)

	staged, _, _, err := g.Status()
	tests.AssertNoError(t, err)

	for _, s := range staged {
		if s != file {
			t.Fatalf("expected `%s` but got `%s`", file, s)
		}
	}
}

func TestCreateTag(t *testing.T) {
	dir := tests.NewTempDir(t)
	defer os.RemoveAll(dir)

	err := os.Chdir(dir)
	tests.AssertNoError(t, err)

	g, err := NewInterface(dir)
	tests.AssertNoError(t, err)

	err = g.Init()
	tests.AssertNoError(t, err)

	file := "test-file.txt"

	err = os.WriteFile(filepath.Join(dir, file), []byte("test"), 0644)
	tests.AssertNoError(t, err)

	err = g.CommitFiles([]string{file}, "Initial commit")
	tests.AssertNoError(t, err)

	tag := "v0.1.0"
	err = g.CreateTag(tag, "Initial release")
	tests.AssertNoError(t, err)

	actual, err := g.LatestTag(true)
	tests.AssertNoError(t, err)

	if tag != actual {
		t.Fatalf("expected `%s` but got `%s`", tag, actual)
	}
}

func TestRemoveFromStaging(t *testing.T) {
	dir := tests.NewTempDir(t)
	defer os.RemoveAll(dir)

	err := os.Chdir(dir)
	tests.AssertNoError(t, err)

	g, err := NewInterface(dir)
	tests.AssertNoError(t, err)

	err = g.Init()
	tests.AssertNoError(t, err)

	file := "test-file.txt"

	err = os.WriteFile(filepath.Join(dir, file), []byte("test"), 0644)
	tests.AssertNoError(t, err)

	err = g.AddToStaging([]string{file})
	tests.AssertNoError(t, err)

	staged, _, _, err := g.Status()
	tests.AssertNoError(t, err)

	for _, s := range staged {
		if s != file {
			t.Fatalf("expected `%s` but got `%s`", file, s)
		}
	}

	err = g.RemoveFromStaging([]string{file}, false)
	tests.AssertNoError(t, err)

	staged, _, _, err = g.Status()
	tests.AssertNoError(t, err)

	if len(staged) > 0 {
		t.Fatalf("expected 0 but got %d", len(staged))
	}
}
