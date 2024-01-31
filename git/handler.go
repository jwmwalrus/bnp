package git

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jwmwalrus/bnp/onerror"
)

// NewHandler retusn a new git interface for the given directory.
func NewHandler(dir string) (Handler, error) {
	if !HasGit() {
		return nil, fmt.Errorf("Unable to find the git command")
	}

	rootDir, err := getRootDir(dir)
	if err != nil {
		// NOTE: not a git directory (yet)
		err = nil
	}

	return &handlerImpl{root: rootDir, log: slog.Default().WithGroup("git")}, nil
}

// handlerImp implements the Handler interface.
type handlerImpl struct {
	root string
	log  *slog.Logger
}

func (h *handlerImpl) AddToStaging(files []string) (err error) {
	files = h.makeAbsPath(files)

	h.log.Info("Staging files", "files", files)

	for _, s := range files {
		if err = h.executeNO("add", s); err != nil {
			_ = h.RemoveFromStaging(files, true)
			return
		}
	}

	return
}

func (h *handlerImpl) Branch() (name string, err error) {
	h.log.Info("Returning active branch")

	out, err := h.execute("branch", "--show-current")
	if err != nil {
		return
	}

	name = strings.TrimSuffix(string(out), "\n")
	return
}

func (h *handlerImpl) Branches(all ...bool) (list []string, err error) {
	h.log.Info("Returning list of branches", "all", all)

	args := []string{"branch", "--no-color"}
	if len(all) > 0 && all[0] {
		args = append(args, "--all")
	}

	out, err := h.execute(args...)
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

func (h *handlerImpl) CheckoutBranch(name string) error {
	h.log.Info("Checing out branch", "name", name)

	return h.executeNO("checkout", name)
}

func (h *handlerImpl) CheckoutNewBranch(name string) error {
	h.log.Info("Checking out new branch", "name", name)

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

func (h *handlerImpl) Commit(msg string) (err error) {
	h.log.Info("Committing", "msg", msg)

	return h.executeNO("commit", "--message", msg)
}

func (h *handlerImpl) CommitFiles(files []string, msg string) (err error) {
	files = h.makeAbsPath(files)

	h.log.With(
		"files", files,
		"msg", msg,
	).Info("Committing with files")

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

func (h *handlerImpl) Config(key string) (value string, err error) {
	h.log.Info("Getting config value", "key", key)

	out, err := h.execute("config", key)
	if err != nil {
		return
	}

	value = strings.TrimSuffix(string(out), "\n")
	return
}

func (h *handlerImpl) DeleteBranch(name string, force ...bool) error {
	h.log.With(
		"name", name,
		"force", force,
	).Info("Deleting branch")

	args := []string{"branch", "--delete"}
	if len(force) > 0 && force[0] {
		args = append(args, "--force")
	}

	args = append(args, name)

	return h.executeNO(args...)
}

func (h *handlerImpl) DiffUpstream(remote string, branch string) (bool, string, error) {
	h.log.With(
		"remote", remote,
		"branch", branch,
	).Info("Checking if current branch differs from upstream's")

	args := []string{"diff", remote + "/" + branch}

	bv, err := h.execute(args...)
	if err != nil {
		return false, "", err
	}

	str := strings.TrimSuffix(string(bv), "\n")

	return len(str) > 0, str, nil
}

func (h *handlerImpl) DropStash(clear ...bool) error {
	args := []string{"stash"}

	h.log.Info("Dropping from stash", "clear", clear)

	if len(clear) > 0 && clear[0] {
		args = append(args, "clear")
	} else {
		args = append(args, "drop")
	}
	return h.executeNO(args...)
}

func (h *handlerImpl) Describe(hash string, exact ...bool) (string, error) {
	h.log.With(
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

	out, err := h.execute(args...)
	if err != nil {
		return "", err
	}

	tag := strings.TrimSuffix(string(out), "\n")
	return tag, nil
}

func (h *handlerImpl) Fetch(remote string) (err error) {
	h.log.Info("Fetching", "remote", remote)

	if remote != "" {
		return h.executeNO("fetch", remote, "--tags")
	}

	return h.executeNO("fetch", "--tags")
}

func (h *handlerImpl) FileChanged(file string) bool {
	files := h.makeAbsPath([]string{file})
	file = files[0]

	h.log.Info("Checking if file changed", "file", file)

	out, err := h.execute("diff", "--name-only", file)
	if err != nil {
		return false
	}
	diff := strings.TrimSuffix(string(out), "\n")
	return len(diff) > 0
}

func (h *handlerImpl) Init(initialBranch string) error {
	h.log.Info("Initializing git repository", "initial-branch", initialBranch)

	if initialBranch != "" {
		return h.executeNO("init", "--initial-branch", initialBranch)
	}

	return h.executeNO("init")
}

func (h *handlerImpl) LatestHash(noFetch ...bool) (hash string, err error) {
	h.log.Info("Getting latest hash", "no-fetch", noFetch)

	doFetch := true
	if len(noFetch) > 0 && noFetch[0] {
		doFetch = false
	}

	if doFetch {
		if err = h.Fetch(""); err != nil {
			return
		}
	}

	out, err := h.execute("rev-parse", "--verify", "HEAD")
	if err != nil {
		return
	}

	hash = strings.TrimSuffix(string(out), "\n")

	return
}

func (h *handlerImpl) LatestTag(noFetch ...bool) (tag string, err error) {
	h.log.Info("Getting latest tag", "no-fetch", noFetch)

	doFetch := true
	if len(noFetch) > 0 && noFetch[0] {
		doFetch = false
	}

	if doFetch {
		if err = h.Fetch(""); err != nil {
			return
		}
	}

	out1, err := h.execute("rev-list", "--tags", "--max-count=1")
	if err != nil {
		return
	}

	hash := strings.TrimSuffix(string(out1), "\n")

	out2, err := h.execute("describe", "--tags", hash)
	if err != nil {
		return
	}

	tag = strings.TrimSuffix(string(out2), "\n")
	return
}

func (h *handlerImpl) Log(maxCount int) ([]LogEntry, error) {
	args := []string{"log", "--pretty=%H;;%at;;%an;;%ae;;%s;;%b"}

	if maxCount > 0 {
		args = append(args, "--max-count", strconv.Itoa(maxCount))
	}

	out, err := h.execute(args...)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")

	list := []LogEntry{}
	for i, l := range lines {
		s := strings.Split(l, ";;")
		if len(s) < 6 {
			slog.Error("Log line contains less elements than expected", "line", i+1)
			continue
		}
		ts, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			slog.With(
				"log-line", i+1,
				"timestamp", s[1],
				"error", err,
			).Error("Failed to parse timestamp")
			return nil, err
		}

		list = append(list, LogEntry{
			Hash:      s[0],
			Timestamp: time.Unix(ts, 0),
			Author:    s[2],
			Email:     s[3],
			Subject:   s[4],
			Body:      s[5],
		})
	}
	return list, nil
}

func (h *handlerImpl) MergeStash(remote, branch, commitMsg string) error {
	h.log.With(
		"remote", remote,
		"branch", branch,
		"commit-msg", commitMsg,
	).Info("Performing Stash + Pull + Merge Stash")

	_, err := h.Stash("", true)
	if err != nil {
		onerror.Log(h.PopStash(""))
		return err
	}

	err = h.Pull(remote, branch, true)
	if err != nil {
		onerror.Log(h.abortMerge())
		onerror.Log(h.PopStash(""))
		return err
	}

	err = h.executeNO("merge", "--squash", "--strategy-option", "theirs", "stash")
	if err != nil {
		onerror.Log(h.PopStash(""))
		return err
	}

	return h.Commit(commitMsg)
}

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

func (h *handlerImpl) NewBranch(name string) error {
	h.log.Info("Creating branch", "name", name)

	return h.executeNO("branch", name)
}

func (h *handlerImpl) NewTag(tag, msg string) (err error) {
	h.log.With(
		"tag", tag,
		"msg", msg,
	).Info("Creating annotated tag")

	return h.executeNO("tag", "--annotate", tag, "-m", msg)
}

func (h *handlerImpl) PopStash(msg string) error {
	h.log.Info("Popping from stash", "stash-msg", msg)

	args := []string{"stash", "pop"}

	if msg != "" {
		list, err := h.StashList()
		if err != nil {
			return err
		}

		for _, l := range list {
			if l.Description == msg {
				args = append(args, l.Name)
				break
			}
		}
	}

	return h.executeNO(args...)
}

func (h *handlerImpl) Pull(remote, branch string, noCommit ...bool) error {
	h.log.With(
		"remote", remote,
		"branch", branch,
		"no-commit", noCommit,
	).Info("Pulling changes from remote branch")

	args := []string{"pull", remote, branch}

	if len(noCommit) > 0 && noCommit[0] {
		args = append(args, "--no-commit")
	}

	return h.executeNO(args...)
}

func (h *handlerImpl) Push(remote, branch string) error {
	h.log.With(
		"remote", remote,
		"branch", branch,
	).Info("Pushing changes to remote branch")

	return h.executeNO("push", remote, branch)
}

func (h *handlerImpl) Remotes() (list map[string]string, err error) {
	h.log.Info("Getting list of remotes")

	out, err := h.execute("remote")
	if err != nil {
		return
	}

	keys := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")

	list = make(map[string]string)
	for _, k := range keys {
		if k == "" {
			continue
		}

		out, err = h.execute("remote", "get-url", k)
		if err != nil {
			return
		}
		list[k] = strings.TrimSuffix(string(out), "\n")
	}

	return
}

func (h *handlerImpl) RemoteUpdate() error {
	h.log.Info("Uodating remote references")

	return h.executeNO("remote", "update")
}

func (h *handlerImpl) RemoveFromStaging(files []string, ignoreErrors ...bool) (err error) {
	files = h.makeAbsPath(files)

	h.log.With(
		"files", files,
		"ignore-errors", ignoreErrors,
	).Info("Removing file(s) from staging")

	ackErrors := true
	if len(ignoreErrors) > 0 && ignoreErrors[0] {
		ackErrors = false
	}

	for _, s := range files {
		if err = h.executeNO("reset", s); err != nil {
			if ackErrors {
				return
			}
		}
	}
	return
}

func (h *handlerImpl) SetConfig(key string, value string) error {
	h.log.With(
		"key", key,
		"value", value,
	).Info("Setting config value")

	return h.executeNO("config", key, value)
}

func (h *handlerImpl) SetRemote(name string, url string) error {
	h.log.With(
		"name", name,
		"url", url,
	).Info("Setting remote")

	list, err := h.Remotes()
	if err != nil {
		return err
	}

	_, ok := list[name]

	if ok {
		return h.executeNO("remote", "set-url", name, url)
	}

	return h.executeNO("remote", "add", name, url)
}

func (h *handlerImpl) SetUpstreamBranchTo(remote, branch string) error {
	h.log.With(
		"remote", remote,
		"branch", branch,
	).Info("Setting active branch to track upsteam's")

	return h.executeNO("branch", "--set-upstream-to", remote+"/"+branch)
}

func (h *handlerImpl) Stash(msg string, untracked ...bool) (entry StashEntry, err error) {
	h.log.With(
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

	err = h.executeNO(args...)
	if err != nil {
		return
	}

	entries, err := h.StashList()
	if err != nil {
		return
	}

	entry = entries[0]

	return
}

func (h *handlerImpl) StashList() ([]StashEntry, error) {
	out, err := h.execute("stash", "list")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(lines) == 0 {
		return nil, err
	}

	list := []StashEntry{}
	for _, l := range lines {
		parts := strings.Split(l, ":")
		if len(parts) == 0 {
			continue
		}

		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		entry := StashEntry{Name: parts[0], Branch: parts[1]}
		if len(parts) > 2 {
			entry.Description = parts[2]
		}
		list = append(list, entry)
	}

	return list, nil
}

func (h *handlerImpl) Status() (staged, unstaged, untracked []string, err error) {
	h.log.Info("Getting status")

	out, err := h.execute("status", "--short", "--porcelain")
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

func (h *handlerImpl) TopLevel() string {
	h.log.Info("Returning top level")

	return h.root
}

func (h *handlerImpl) Unstage(files []string) error {
	h.log.Info("Removing files from staging", "files", files)

	args := []string{"reset"}
	args = append(args, files...)

	return h.executeNO(args...)
}

func (h *handlerImpl) abortMerge() error {
	h.log.Info("Aborting merge")

	return h.executeNO("merge", "--abort")
}

func (h *handlerImpl) executeNO(in ...string) error {
	args := []string{"-C", h.root}
	args = append(args, in...)

	cmd := exec.Command("git", args...)

	slog.Debug("Running git command with combined output", "cmd", cmd)

	_, err := cmd.CombinedOutput()
	return err
}

func (h *handlerImpl) execute(in ...string) ([]byte, error) {
	args := []string{"-C", h.root}
	args = append(args, in...)

	cmd := exec.Command("git", args...)
	out := &bytes.Buffer{}
	cmd.Stdout = out

	slog.Debug("Running git command", "cmd", cmd)

	err := cmd.Run()
	return out.Bytes(), err
}

func (h *handlerImpl) makeAbsPath(files []string) []string {
	pwd, err := os.Getwd()
	if err != nil || (pwd != h.root && !strings.Contains(pwd, h.root)) {
		h.log.With(
			"pwd", pwd,
			"root", h.root,
			"error", err,
		).Warn("Failed check for absolute path")
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

	out, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	if err != nil {
		rootDir = dir
		return
	}

	rootDir = strings.TrimSuffix(string(out), "\n")
	return
}
