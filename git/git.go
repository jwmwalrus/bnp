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

// Interface provides a handler to git's command line
type Interface interface {
	// AddToStaging adds the given files to staging
	AddToStaging(files []string) (err error)

	// CreateTag creates an annotated tag
	CreateTag(tag, msg string) (err error)

	// Commit commits files in staging with the given message
	Commit(msg string) (err error)

	// CommitFiles adds the given files to staging and commits them with the given message
	CommitFiles(files []string, msg string) (err error)

	// Describe returns the corresponding tag for the given hash
	Describe(hash string, exact ...bool) (tag string, err error)

	// Fetch brings the latest changes for the given remote
	Fetch(remote string) (err error)

	// FileChanged checks if a file changed and should be added to staging
	FileChanged(file string) bool

	// Init git-initializes the root directory
	Init() error

	// LatestTag Returns the latest tag for the git repo related to the working directory
	LatestTag(noFetch bool) (tag string, err error)

	// MoveToRootDir changes working directory to git's root
	MoveToRootDir() RestoreCwdFunc

	// RemoveFromStaging removes the given files from the stagin area
	RemoveFromStaging(files []string, ignoreErrors bool) (err error)

	// Status reports the current status of the working tree
	Status() (staged, unstaged, untracked []string, err error)
}

func NewInterface(dir string) (Interface, error) {
	if !HasGit() {
		return nil, fmt.Errorf("Unable to find the git command")
	}

	rootDir, err := getRootDir(dir)
	if err != nil {
		// NOTE: not a git directory (yet)
		err = nil
	}

	return &handler{root: rootDir}, nil
}

type handler struct {
	root string
}

// AddToStaging adds the given files to staging
func (h *handler) AddToStaging(files []string) (err error) {
	for i := range files {
		files[i], err = filepath.Abs(files[i])
		if err != nil {
			return
		}
	}

	restore := h.MoveToRootDir()
	defer restore()

	fmt.Printf("\tStaging files...\n")
	for _, s := range files {
		if err = executeNO("add", s); err != nil {
			_ = h.RemoveFromStaging(files, true)
			return
		}
	}

	return
}

// CreateTag creates an annotated tag
func (h *handler) CreateTag(tag, msg string) (err error) {
	restore := h.MoveToRootDir()
	defer restore()

	fmt.Printf("\nCreating annotated tag %v with message '%v'...\n", tag, msg)

	return executeNO("tag", "--annotate", tag, "-m", msg)
}

// Commit commits files in staging with the given message
func (h *handler) Commit(msg string) (err error) {
	restore := h.MoveToRootDir()
	defer restore()

	fmt.Printf("\tCommiting with '%v' as message...\n", msg)
	return executeNO("commit", "--message", msg)
}

// CommitFiles adds the given files to staging and commits them with the given message
func (h *handler) CommitFiles(files []string, msg string) (err error) {
	for i := range files {
		files[i], err = filepath.Abs(files[i])
		if err != nil {
			return
		}
	}

	restore := h.MoveToRootDir()
	defer restore()

	fmt.Printf("\nCommiting files...\n")

	if err = h.AddToStaging(files); err != nil {
		_ = h.RemoveFromStaging(files, true)
		return
	}

	if err = h.Commit(msg); err != nil {
		_ = h.RemoveFromStaging(files, true)
		return
	}

	return
}

// Describe returns the corresponding tag for the given hash
func (h *handler) Describe(hash string, exact ...bool) (string, error) {
	restore := h.MoveToRootDir()
	defer restore()

	args := []string{"describe"}
	if len(exact) > 0 && exact[0] {
		args = append(args, "--match-exact")
	}
	args = append(args, hash)

	out, err := execute(args...)
	if err != nil {
		return "", err
	}

	tag := string(out)
	tag = strings.TrimSuffix(tag, "\n")
	return tag, nil
}

// Fetch brings the latest changes for the given remote
func (h *handler) Fetch(remote string) (err error) {
	fmt.Printf("\tFetching...\n")

	if remote != "" {
		return executeNO("fetch", remote, "--tags")
	}

	return executeNO("fetch", "--tags")
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

	out, err := execute("diff", "--name-only", file)
	if err != nil {
		return false
	}
	diff := string(out)
	diff = strings.TrimSuffix(diff, "\n")
	return len(diff) > 0
}

// Init git-initializes the root directory
func (h *handler) Init() error {
	restore := h.MoveToRootDir()
	defer restore()

	if _, _, _, err := h.Status(); err == nil {
		return nil
	}

	return executeNO("init")
}

// LatestTag Returns the latest tag for the git repo related to the working directory
func (h *handler) LatestTag(noFetch bool) (tag string, err error) {
	restore := h.MoveToRootDir()
	defer restore()

	fmt.Printf("\nGetting latest git tag...\n")

	if !noFetch {
		fmt.Printf("\tFetching...\n")
		if err = h.Fetch(""); err != nil {
			fmt.Printf("...fetching failed!\n")
		}
	}

	out1, err := execute("rev-list", "--tags", "--max-count=1")
	if err != nil {
		return
	}

	hash := string(out1)
	hash = strings.TrimSuffix(hash, "\n")

	out2, err := execute("describe", "--tags", hash)
	if err != nil {
		return
	}

	tag = string(out2)
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
func (h *handler) RemoveFromStaging(files []string, ignoreErrors bool) (err error) {
	for i := range files {
		files[i], err = filepath.Abs(files[i])
		if err != nil {
			return
		}
	}

	restore := h.MoveToRootDir()
	defer restore()

	for _, s := range files {
		if err = executeNO("reset", s); err != nil {
			if !ignoreErrors {
				return
			}
		}
	}
	return
}

// Status reports the current status of the working tree
func (h *handler) Status() (staged, unstaged, untracked []string, err error) {
	out, err := execute("status", "--short", "--porcelain")
	if err != nil {
		return
	}

	all := strings.Split(string(out), "\n")

	for _, line := range all {
		if len(line) < 3 {
			continue
		}

		state := line[:2]
		file := line[3:]

		switch state[0] {
		case '?', '!':
			before, _, _ := strings.Cut(file, "->")
			untracked = append(untracked, before)
		case ' ':
			switch state[1] {
			case 'A', 'C', 'D', 'M', 'R', 'T', 'U':
				before, _, _ := strings.Cut(file, "->")
				unstaged = append(unstaged, before)
			default:
			}
		case 'A', 'C', 'D', 'M', 'R', 'T', 'U':
			before, _, _ := strings.Cut(file, "->")
			staged = append(staged, before)
		default:
		}
	}

	return
}

func (h *handler) TopLevel(dir string) (rootDir string, err error) {
	rootDir, err = getRootDir(dir)
	if err != nil {
		rootDir = ""
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

	if pwd != dir {
		err = os.Chdir(dir)
		if err != nil {
			return
		}

		defer func() { os.Chdir(pwd) }()
	}

	out, err := execute("rev-parse", "--show-toplevel")
	if err != nil {
		rootDir = dir
		return
	}

	rootDir = string(out)
	rootDir = strings.TrimSuffix(rootDir, "\n")
	return
}

func executeNO(args ...string) error {
	_, err := exec.Command("git", args...).CombinedOutput()
	return err
}

func execute(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	err := cmd.Run()
	return out.Bytes(), err
}
