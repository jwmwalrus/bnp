package git

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RestoreCwdFunc defines the signature of the closure to restore the working directory
type RestoreCwdFunc func() error

// CreateTag creates an annotated tag
func CreateTag(tag, msg string) (err error) {
	if !HasGit() {
		err = errors.New("Unable to find the git command")
		return
	}
	fmt.Printf("\nCreating annotated tag %v with message '%v'...\n", tag, msg)

	_, err = exec.Command("git", "tag", "-a", tag, "-m", msg).CombinedOutput()
	return
}

// CommitFiles adds the given files to staging and commits them with the given message
func CommitFiles(sList []string, m string) (err error) {
	if !HasGit() {
		err = errors.New("Unable to find the git command")
		return
	}
	fmt.Printf("\nCommiting files...\n")

	fmt.Printf("\tStaging files...\n")
	for _, s := range sList {
		if _, err = exec.Command("git", "add", s).CombinedOutput(); err != nil {
			_ = RemoveFromStaging(sList, true)
			return
		}
	}
	fmt.Printf("\tCommiting with '%v' as message...\n", m)
	if _, err = exec.Command("git", "commit", "-m", m).CombinedOutput(); err != nil {
		_ = RemoveFromStaging(sList, true)
		return
	}

	return
}

// Describe returns the corresponding tag for the given hash
func Describe(hash string, exact ...bool) (string, error) {
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
func FileChanged(file string) bool {
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
func GetLatestTag(noFetch bool) (tag string, err error) {
	if !HasGit() {
		err = errors.New("Unable to find the git command")
		return
	}
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

// HasGit checks if the git command exists in PATH
func HasGit() bool {
	s, err := exec.LookPath("git")
	return s != "" && err == nil
}

// MoveToRootDir changes working directory to git's root
func MoveToRootDir() (RestoreCwdFunc, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	root := pwd

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output := &bytes.Buffer{}
	cmd.Stdout = output
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	root = string(output.Bytes())
	root = strings.TrimSuffix(root, "\n")
	if root == pwd {
		return func() error { return nil }, nil
	}

	if err = os.Chdir(root); err != nil {
		return nil, err
	}

	return func() error { return os.Chdir(pwd) }, nil
}

// RemoveFromStaging removes the given files from the stagin area
func RemoveFromStaging(sList []string, ignoreErrors bool) (err error) {
	if !HasGit() {
		err = errors.New("Unable to find the git command")
		return
	}

	for _, s := range sList {
		if _, err = exec.Command("git", "reset", s).CombinedOutput(); err != nil {
			if !ignoreErrors {
				return
			}
		}
	}
	return
}
