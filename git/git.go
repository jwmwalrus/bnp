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

// Handler provides a handler to git's command line
type Handler interface {
	// AddToStaging adds the given files to staging
	AddToStaging(files []string) (err error)

	// Branch returns the active branch
	Branch() (name string, err error)

	// Branches returns the list of branches
	Branches(all ...bool) (list []string, err error)

	// CheckoutBranch checks out the given branch
	CheckoutBranch(name string) error

	// CheckoutNewBranch creates the given branch and checks it out
	CheckoutNewBranch(name string) error

	// Commit commits files in staging with the given message
	Commit(msg string) (err error)

	// CommitFiles adds the given files to staging and commits them with the given message
	CommitFiles(files []string, msg string) (err error)

	// Config returns the current config value
	Config(key string) (value string, err error)

	// DeleteBranch removes the given branch
	DeleteBranch(name string, force ...bool) error

	// Describe returns the corresponding tag for the given hash
	Describe(hash string, exact ...bool) (tag string, err error)

	// DropStash
	DropStash(all ...bool) error

	// Fetch brings the latest changes for the given remote
	Fetch(remote string) (err error)

	// FileChanged checks if a file changed and should be added to staging
	FileChanged(file string) bool

	// Init git-initializes the root directory
	Init(initialBranch string) error

	// LatestHash Returns the latest tag for the git repo related to the working directory
	LatestHash(noFetch ...bool) (hash string, err error)

	// LatestTag Returns the latest tag for the git repo related to the working directory
	LatestTag(noFetch ...bool) (tag string, err error)

	// MergeStash merges remote changes, preserving ours
	MergeStash(remote, branch, commitMsg string) error

	// MustMoveToRootDir changes working directory to git's root
	MustMoveToRootDir() RestoreCwdFunc

	// NewBranch creates a new branch
	NewBranch(name string) error

	// NewTag creates an annotated tag
	NewTag(tag, msg string) (err error)

	// Pull updates tree with remote changes
	Pull(remote, branch string) error

	// Push sends branch changes to remote
	Push(remote, branch string) error

	// Remotes returns the list of remotes set for the repository
	Remotes() (list map[string]string, err error)

	// RemoveFromStaging removes the given files from the stagin area
	RemoveFromStaging(files []string, ignoreErrors ...bool) (err error)

	// SetUpstreamBranchTo implements the Handler interface
	SetUpstreamBranchTo(remote, branch string) error

	// SetConfig sets a config value
	SetConfig(key, value string) error

	// SetRemote adds remote or sets URL for an existing remote
	SetRemote(name, url string) error

	// Stash stashes local changes
	Stash(msg string, untracked ...bool) (stash string, err error)

	// Status reports the current status of the working tree
	Status() (staged, unstaged, untracked []string, err error)

	// TopLevel returns the root directory
	TopLevel(dir string) (rootDir string, err error)
}

// NewHandler retusn a new git interface for the given directory
func NewHandler(dir string) (Handler, error) {
	if !HasGit() {
		return nil, fmt.Errorf("Unable to find the git command")
	}

	rootDir, err := getRootDir(dir)
	if err != nil {
		// NOTE: not a git directory (yet)
		err = nil
	}

	return &handlerImpl{root: rootDir}, nil
}

type handlerImpl struct {
	root string
}

// AddToStaging implements the Handler interface
func (h *handlerImpl) AddToStaging(files []string) (err error) {
	files = h.makeAbsPath(files)

	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Staging files", "files", files)

	for _, s := range files {
		if err = executeNO("add", s); err != nil {
			_ = h.RemoveFromStaging(files, true)
			return
		}
	}

	return
}

// Branch implements the Handler interface
func (h *handlerImpl) Branch() (name string, err error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Returning active branch")

	out, err := execute("branch", "--show-current")
	if err != nil {
		return
	}

	name = strings.TrimSuffix(string(out), "\n")
	return
}

// Branches implements the Handler interface
func (h *handlerImpl) Branches(all ...bool) (list []string, err error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Returning list of branches", "all", all)

	args := []string{"branch", "--no-color"}
	if len(all) > 0 && all[0] {
		args = append(args, "--all")
	}

	out, err := execute(args...)
	if err != nil {
		return
	}

	list = strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	for i := range list {
		list[i] = strings.TrimPrefix(list[i], "*")
		list[i] = strings.TrimSpace(list[i])
	}

	return
}

// CheckoutBranch implements the Handler interface
func (h *handlerImpl) CheckoutBranch(name string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Checing out branch", "name", name)

	return executeNO("checkout", name)
}

// CheckoutNewBranch implements the Handler interface
func (h *handlerImpl) CheckoutNewBranch(name string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Checking out new branch", "name", name)

	err := h.NewBranch(name)
	if err != nil {
		return err
	}

	err = h.CheckoutBranch(name)
	if err != nil {
		h.DeleteBranch(name, false)
		return err
	}

	return nil
}

// Commit implements the Handler interface
func (h *handlerImpl) Commit(msg string) (err error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Commiting", "msg", msg)

	return executeNO("commit", "--message", msg)
}

// CommitFiles implements the Handler interface
func (h *handlerImpl) CommitFiles(files []string, msg string) (err error) {
	files = h.makeAbsPath(files)

	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"files", files,
		"msg", msg,
	).Info("Commiting with files")

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

// Config implements the Handler interface
func (h *handlerImpl) Config(key string) (value string, err error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Getting config value", "key", key)

	out, err := execute("config", key)
	if err != nil {
		return
	}

	value = strings.TrimSuffix(string(out), "\n")
	return
}

// DeleteBranch implements the Handler interface
func (h *handlerImpl) DeleteBranch(name string, force ...bool) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"name", name,
		"force", force,
	).Info("Deleting branch")

	args := []string{"branch", "--delete"}
	if len(force) > 0 && force[0] {
		args = append(args, "--force")
	}

	args = append(args, name)

	return executeNO(args...)
}

func (h *handlerImpl) DropStash(clear ...bool) error {
	args := []string{"stash"}

	slog.Info("Dropping from stash", "clear", clear)

	if len(clear) > 0 && clear[0] {
		args = append(args, "clear")
	} else {
		args = append(args, "drop")
	}
	return executeNO(args...)
}

// Describe implements the Handler interface
func (h *handlerImpl) Describe(hash string, exact ...bool) (string, error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"hash", hash,
		"exact", exact,
	).Info("Describing current tree state")

	args := []string{"describe"}
	if len(exact) > 0 && exact[0] {
		args = append(args, "--match-exact")
	}
	if len(hash) > 0 {
		args = append(args, hash)
	}

	out, err := execute(args...)
	if err != nil {
		return "", err
	}

	tag := string(out)
	tag = strings.TrimSuffix(tag, "\n")
	return tag, nil
}

// Fetch implements the Handler interface
func (h *handlerImpl) Fetch(remote string) (err error) {
	slog.Info("Fetching", "remote", remote)

	if remote != "" {
		return executeNO("fetch", remote, "--tags")
	}

	return executeNO("fetch", "--tags")
}

// FileChanged implements the Handler interface
func (h *handlerImpl) FileChanged(file string) bool {
	files := h.makeAbsPath([]string{file})
	file = files[0]

	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Checking if file changed", "file", file)

	out, err := execute("diff", "--name-only", file)
	if err != nil {
		return false
	}
	diff := string(out)
	diff = strings.TrimSuffix(diff, "\n")
	return len(diff) > 0
}

// Init implements the Handler interface
func (h *handlerImpl) Init(initialBranch string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Initializing git repository", "initial-branch", initialBranch)

	if _, _, _, err := h.Status(); err == nil {
		return nil
	}

	if initialBranch != "" {
		return executeNO("init", "--initial-branch", initialBranch)
	}

	return executeNO("init")
}

// LatestHash implements the Handler interface
func (h *handlerImpl) LatestHash(noFetch ...bool) (hash string, err error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Getting latest hash", "no-fetch", noFetch)

	doFetch := true
	if len(noFetch) > 0 && noFetch[0] {
		doFetch = false
	}

	if doFetch {
		if err = h.Fetch(""); err != nil {
			return
		}
	}

	out1, err := execute("rev-parse", "--verify", "HEAD")
	if err != nil {
		return
	}

	hash = string(out1)
	hash = strings.TrimSuffix(hash, "\n")

	return
}

// LatestTag implements the Handler interface
func (h *handlerImpl) LatestTag(noFetch ...bool) (tag string, err error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Getting latest tag", "no-fetch", noFetch)

	doFetch := true
	if len(noFetch) > 0 && noFetch[0] {
		doFetch = false
	}

	if doFetch {
		if err = h.Fetch(""); err != nil {
			return
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

// MergeStash implements the Handler interface
func (h *handlerImpl) MergeStash(remote, branch, commitMsg string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"remote", remote,
		"branch", branch,
		"commit-msg", commitMsg,
	).Info("Performing Stash + Pull + Merge Stash")

	_, err := h.Stash("", true)
	if err != nil {
		return err
	}

	err = h.Pull(remote, branch)
	if err != nil {
		return err
	}

	err = executeNO("merge", "--squash", "--strategy-option", "theirs", "stash")
	if err != nil {
		return err
	}

	return h.Commit(commitMsg)
}

// MustMoveToRootDir implements the Handler interface
func (h *handlerImpl) MustMoveToRootDir() RestoreCwdFunc {
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

// NewBranch implements the Handler interface
func (h *handlerImpl) NewBranch(name string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Creating branch", "name", name)

	return executeNO("branch", name)
}

// NewTag implements the Handler interface
func (h *handlerImpl) NewTag(tag, msg string) (err error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"tag", tag,
		"msg", msg,
	).Info("Creating annotated tag")

	return executeNO("tag", "--annotate", tag, "-m", msg)
}

// Pull implements the Handler interface
func (h *handlerImpl) Pull(remote, branch string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"remote", remote,
		"branch", branch,
	).Info("Pulling changes from remote branch")

	return executeNO("pull", remote, branch)
}

// Push implements the Handler interface
func (h *handlerImpl) Push(remote, branch string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"remote", remote,
		"branch", branch,
	).Info("Pushing changes to remote branch")

	return executeNO("push", remote, branch)
}

// Remotes implements the Handler interface
func (h *handlerImpl) Remotes() (list map[string]string, err error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.Info("Getting list of remotes")

	out, err := execute("remote")
	if err != nil {
		return
	}

	keys := strings.Split(string(out), "\n")

	list = make(map[string]string)
	for _, k := range keys {
		if k == "" {
			continue
		}

		out, err = execute("remote", "get-url", k)
		if err != nil {
			return
		}
		list[k] = string(out)
	}

	return
}

// RemoveFromStaging implements the Handler interface
func (h *handlerImpl) RemoveFromStaging(files []string, ignoreErrors ...bool) (err error) {
	files = h.makeAbsPath(files)

	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"files", files,
		"ignore-errors", ignoreErrors,
	).Info("Removing file(s) from staging")

	ackErrors := true
	if len(ignoreErrors) > 0 && ignoreErrors[0] {
		ackErrors = false
	}

	for _, s := range files {
		if err = executeNO("reset", s); err != nil {
			if ackErrors {
				return
			}
		}
	}
	return
}

// SetConfig implements the Handler interface
func (h *handlerImpl) SetConfig(key string, value string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"key", key,
		"value", value,
	).Info("Setting config value")

	return executeNO("config", key, value)
}

// SetRemote implements the Handler interface
func (h *handlerImpl) SetRemote(name string, url string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"name", name,
		"url", url,
	).Info("Setting remote")

	list, err := h.Remotes()
	if err != nil {
		return err
	}

	_, ok := list[name]

	if ok {
		return executeNO("remote", "set-url", name, url)
	}

	return executeNO("remote", "add", name, url)
}

// SetUpstreamBranchTo implements the Handler interface
func (h *handlerImpl) SetUpstreamBranchTo(remote, branch string) error {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"remote", remote,
		"branch", branch,
	).Info("Setting active branch to track upsteam's")

	return executeNO("branch", "--set-upstream-to", remote+"/"+branch)
}

// Stash implements the handler interface
func (h *handlerImpl) Stash(msg string, untracked ...bool) (stash string, err error) {
	restore := h.MustMoveToRootDir()
	defer restore()

	slog.With(
		"msg", msg,
		"untracked", untracked,
	).Info("Stashing changes")

	args := []string{"stash"}

	if msg != "" {
		args = append(args, "--message", msg)
	}

	if len(untracked) > 0 && untracked[0] {
		args = append(args, "--include-untracked")
	}

	err = executeNO(args...)

	out, err := execute("stash", "list")
	if err != nil {
		return
	}

	list := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(list) == 0 {
		return
	}

	parts := strings.Split(list[0], ":")
	if len(parts) == 0 {
		return
	}

	stash = parts[0]

	return
}

// Status implements the Handler interface
func (h *handlerImpl) Status() (staged, unstaged, untracked []string, err error) {
	slog.Info("Getting status")

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

// TopLevel implements the Handler interface
func (h *handlerImpl) TopLevel(dir string) (rootDir string, err error) {
	slog.Info("Getting top level", "dir", dir)

	rootDir, err = getRootDir(dir)
	if err != nil {
		rootDir = ""
	}
	return
}

func (h *handlerImpl) makeAbsPath(files []string) []string {
	pwd, err := os.Getwd()
	if err != nil || (pwd != h.root && !strings.Contains(pwd, h.root)) {
		slog.With(
			"pwd", pwd,
			"root", h.root,
			"error", err,
		).Error("Failed check for absolute path")
		return files
	}

	var newFiles []string
	for i := range files {
		f, err := filepath.Abs(files[i])
		if err != nil {
			return files
		}
		newFiles = append(newFiles, f)
	}

	return newFiles
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
	cmd := exec.Command("git", args...)
	slog.Debug("Running command with combined output", "cmd", cmd)
	_, err := cmd.CombinedOutput()
	return err
}

func execute(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	slog.Debug("Running command", "cmd", cmd)
	err := cmd.Run()
	return out.Bytes(), err
}
