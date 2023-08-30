package git

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jwmwalrus/bnp/onerror"
)

// RestoreCwdFunc defines the signature of the closure to restore the working directory
type RestoreCwdFunc func() error

type Interface interface {
	// CreateTag creates an annotated tag
	CreateTag(tag, msg string) (err error)

	// CommitFiles adds the given files to staging and commits them with the given message
	CommitFiles(sList []string, msg string) (err error)

	// Describe returns the corresponding tag for the given hash
	Describe(hash string, exact ...bool) (string, error)

	// FileChanged checks if a file changed and should be added to staging
	FileChanged(file string) bool

	// GetLatestTag Returns the latest tag for the git repo related to the working directory
	GetLatestTag(noFetch bool) (tag string, err error)

	// MoveToRootDir changes working directory to git's root
	MoveToRootDir() RestoreCwdFunc

	// RemoveFromStaging removes the given files from the stagin area
	RemoveFromStaging(sList []string, ignoreErrors bool) (err error)
}

func NewInterface(dir string) (Interface, error) {
	if !HasGit() {
		return nil, fmt.Errorf("Unable to find the git command")
	}

	rootDir, err := getRootDir(dir)
	if err != nil {
		return nil, err
	}

	return &handler{root: rootDir}, nil
}

type handler struct {
	root string
}

// CreateTag creates an annotated tag
func (h *handler) CreateTag(tag, msg string) (err error) {
	restore := h.MoveToRootDir()
	defer restore()

	fmt.Printf("\nCreating annotated tag %v with message '%v'...\n", tag, msg)

	_, err = exec.Command("git", "tag", "-a", tag, "-m", msg).CombinedOutput()
	return
}

// CommitFiles adds the given files to staging and commits them with the given message
func (h *handler) CommitFiles(sList []string, msg string) (err error) {
	for i := range sList {
		sList[i], err = filepath.Abs(sList[i])
		if err != nil {
			return
		}
	}

	restore := h.MoveToRootDir()
	defer restore()

	fmt.Printf("\nCommiting files...\n")

	fmt.Printf("\tStaging files...\n")
	for _, s := range sList {
		if _, err = exec.Command("git", "add", s).CombinedOutput(); err != nil {
			_ = h.RemoveFromStaging(sList, true)
			return
		}
	}
	fmt.Printf("\tCommiting with '%v' as message...\n", msg)
	if _, err = exec.Command("git", "commit", "-m", msg).CombinedOutput(); err != nil {
		_ = h.RemoveFromStaging(sList, true)
		return
	}

	return
}

// Describe returns the corresponding tag for the given hash
func (h *handler) Describe(hash string, exact ...bool) (string, error) {
	restore := h.MoveToRootDir()
	defer restore()

	list := []string{"describe"}
	if len(exact) > 0 && exact[0] {
		list = append(list, "--match-exact")
	}
	list = append(list, hash)
	cmd1 := exec.Command("git", list...)
	output1 := &bytes.Buffer{}
	cmd1.Stdout = output1
	if err := cmd1.Run(); err != nil {
		return "", err
	}
	tag := string(output1.Bytes())
	tag = strings.TrimSuffix(tag, "\n")
	return tag, nil
}

// FileChanged checks if a file changed and should be added to staging
func (h *handler) FileChanged(file string) bool {
	file, err := filepath.Abs(file)
	if err != nil {
		slog.With(
			"file", file,
			"error", err,
		).Error("Failed to get file's absolute path")
		return false
	}

	restore := h.MoveToRootDir()
	defer restore()

	cmd1 := exec.Command("git", "diff", "--name-only", file)
	output1 := &bytes.Buffer{}
	cmd1.Stdout = output1
	if err := cmd1.Run(); err != nil {
		return false
	}
	diff := string(output1.Bytes())
	diff = strings.TrimSuffix(diff, "\n")
	return len(diff) > 0
}

// GetLatestTag Returns the latest tag for the git repo related to the working directory
func (h *handler) GetLatestTag(noFetch bool) (tag string, err error) {
	restore := h.MoveToRootDir()
	defer restore()

	fmt.Printf("\nGetting latest git tag...\n")

	if !noFetch {
		fmt.Printf("\tFetching...\n")
		if _, err = exec.Command("git", "fetch", "--tags").CombinedOutput(); err != nil {
			fmt.Printf("...fetching failed!\n")
		}
	}

	cmd1 := exec.Command("git", "rev-list", "--tags", "--max-count=1")
	output1 := &bytes.Buffer{}
	cmd1.Stdout = output1
	if err = cmd1.Run(); err != nil {
		return
	}
	hash := string(output1.Bytes())
	hash = strings.TrimSuffix(hash, "\n")

	cmd2 := exec.Command("git", "describe", "--tags", hash)
	output2 := &bytes.Buffer{}
	cmd2.Stdout = output2
	if err = cmd2.Run(); err != nil {
		return
	}

	tag = string(output2.Bytes())
	tag = strings.TrimSuffix(tag, "\n")

	return
}

// MoveToRootDir changes working directory to git's root
func (h *handler) MoveToRootDir() RestoreCwdFunc {
	cwd, err := os.Getwd()
	onerror.Fatal(err)
	root := cwd

	if root == h.root {
		return func() error { return nil }
	}

	root = h.root
	onerror.Fatal(os.Chdir(root))

	return func() error { return os.Chdir(cwd) }
}

// RemoveFromStaging removes the given files from the stagin area
func (h *handler) RemoveFromStaging(sList []string, ignoreErrors bool) (err error) {
	for i := range sList {
		sList[i], err = filepath.Abs(sList[i])
		if err != nil {
			return
		}
	}

	restore := h.MoveToRootDir()
	defer restore()

	for _, s := range sList {
		if _, err = exec.Command("git", "reset", s).CombinedOutput(); err != nil {
			if !ignoreErrors {
				return
			}
		}
	}
	return
}

// HasGit checks if the git command exists in PATH
func HasGit() bool {
	s, err := exec.LookPath("git")
	return s != "" && err == nil
}

func getRootDir(dir string) (rootDir string, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	err = os.Chdir(dir)
	if err != nil {
		return
	}

	defer func() { os.Chdir(pwd) }()

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output := &bytes.Buffer{}
	cmd.Stdout = output
	err = cmd.Run()
	if err != nil {
		return
	}

	rootDir = string(output.Bytes())
	rootDir = strings.TrimSuffix(rootDir, "\n")

	return
}
