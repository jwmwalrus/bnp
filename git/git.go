package git

import (
	"os/exec"
	"time"
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

	// Log Returns log entries
	Log(maxCount int) ([]LogEntry, error)

	// MergeStash merges remote changes, preserving ours
	MergeStash(remote, branch, commitMsg string) error

	// MustMoveToRootDir changes working directory to git's root
	MustMoveToRootDir() RestoreCwdFunc

	// NewBranch creates a new branch
	NewBranch(name string) error

	// NewTag creates an annotated tag
	NewTag(tag, msg string) (err error)

	// PopStash pops the most recent stash
	PopStash(msg string) error

	// Pull updates tree with remote changes
	Pull(remote, branch string, noCommit ...bool) error

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
	Stash(msg string, untracked ...bool) (StashEntry, error)

	// StashList returns the list of stash entries
	StashList() ([]StashEntry, error)

	// Status reports the current status of the working tree
	Status() (staged, unstaged, untracked []string, err error)

	// TopLevel returns the root directory
	TopLevel() string
}

type LogEntry struct {
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
}

type StashEntry struct {
	Name        string `json:"name"`
	Branch      string `json:"branch"`
	Description string `json:"description"`
}

// HasGit checks if the git command exists in PATH
func HasGit() bool {
	s, err := exec.LookPath("git")
	return s != "" && err == nil
}
