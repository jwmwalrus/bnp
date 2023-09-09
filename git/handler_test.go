package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jwmwalrus/bnp/tests"
	"github.com/jwmwalrus/bnp/tests/assert"
)

// MergeStash

func TestAddToStaging(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	file := "test-file.txt"

	err := os.WriteFile(filepath.Join(dir, file), []byte("test"), 0644)
	assert.NoError(t, err)

	err = g.AddToStaging([]string{file})
	assert.NoError(t, err)

	staged, _, _, err := g.Status()
	assert.NoError(t, err)

	assert.Equal(t, []string{file}, staged)
}

func TestBranch(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	current, err := g.Branch()
	assert.NoError(t, err)

	if current != "main" {
		t.Fatalf("expected `main` but got `%s`", current)
	}
}

func TestBranches(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	file := "file.txt"
	err := os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
	assert.NoError(t, err)

	err = g.CommitFiles([]string{file}, "Initial commit")
	assert.NoError(t, err)

	err = g.CheckoutNewBranch("test")
	assert.NoError(t, err)

	list, err := g.Branches()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(list))

	assert.Equal(t, []string{"main", "test"}, list)
}

func TestCheckoutNewBranch(t *testing.T) {
	testCases := []struct {
		name    string
		commit  bool
		wantErr bool
	}{
		{
			name:   "success",
			commit: true,
		},
		{
			name:    "branch not created",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, dir := newTestRepo(t, "main")
			defer os.RemoveAll(dir)

			if tc.commit {
				file := "file.txt"
				err := os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
				assert.NoError(t, err)

				err = g.CommitFiles([]string{file}, "Initial commit")
				assert.NoError(t, err)
			}

			err := g.CheckoutNewBranch("test")
			assert.Equal(t, tc.wantErr, err != nil)

			if tc.wantErr {
				return
			}

			current, err := g.Branch()
			assert.NoError(t, err)

			assert.Equal(t, "test", current)
		})
	}
}

func TestDeleteBranch(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	file := "file.txt"
	err := os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
	assert.NoError(t, err)

	err = g.CommitFiles([]string{file}, "Initial commit")
	assert.NoError(t, err)

	err = g.NewBranch("test")
	assert.NoError(t, err)

	list, err := g.Branches()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(list))
	assert.Equal(t, []string{"main", "test"}, list)

	err = g.DeleteBranch("test")
	assert.NoError(t, err)

	list, err = g.Branches()
	assert.NoError(t, err)

	assert.Equal(t, []string{"main"}, list)
}

func TestDescribe(t *testing.T) {
	testCases := []struct {
		name    string
		commit  bool
		tag     bool
		wantErr bool
	}{
		{
			name:   "success",
			commit: true,
			tag:    true,
		},
		{
			name:    "no commits",
			commit:  true,
			wantErr: true,
		},
		{
			name:    "no tag",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, dir := newTestRepo(t, "main")
			defer os.RemoveAll(dir)

			if tc.commit {
				file := "file.txt"
				err := os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
				assert.NoError(t, err)

				err = g.CommitFiles([]string{file}, "Initial commit")
				assert.NoError(t, err)
			}

			tag := "v0.1.0"
			if tc.tag {
				err := g.NewTag(tag, "Initial release")
				assert.NoError(t, err)
			}

			actual, err := g.Describe("")
			assert.Equal(t, tc.wantErr, err != nil)

			if tc.wantErr {
				return
			}

			assert.Equal(t, tag, actual)
		})
	}
}

func TestFetch(t *testing.T) {
	remote := tests.NewTempDir(t)
	defer os.RemoveAll(remote)

	_, err := exec.Command("git", "init", "--bare", remote).CombinedOutput()
	assert.NoError(t, err)

	g1, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	err = g1.SetRemote("origin", remote)
	assert.NoError(t, err)

	file := "file.txt"
	err = os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
	assert.NoError(t, err)

	err = g1.CommitFiles([]string{file}, "Initial commit")
	assert.NoError(t, err)

	err = g1.Push("origin", "main")
	assert.NoError(t, err)

	hash, err := g1.LatestHash(true)
	assert.NoError(t, err)

	g2, dir2 := newTestRepo(t, "")
	defer os.RemoveAll(dir2)

	err = g2.SetRemote("origin", remote)
	assert.NoError(t, err)

	err = g2.Fetch("origin")
	assert.NoError(t, err)

	err = g2.CheckoutBranch("main")
	assert.NoError(t, err)

	actual, err := g2.LatestHash(true)
	assert.NoError(t, err)

	assert.Equal(t, hash, actual)
}

func TestFileChanged(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	file := "file.txt"
	err := os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
	assert.NoError(t, err)

	err = g.CommitFiles([]string{file}, "Initial commit")
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(dir, file), []byte("test"), 0644)
	assert.NoError(t, err)

	changed := g.FileChanged(file)
	assert.Equal(t, true, changed)

	_, unstaged, _, err := g.Status()
	assert.Equal(t, true, changed)

	assert.Equal(t, []string{file}, unstaged)
}

func TestMergeStash(t *testing.T) {
	remote := tests.NewTempDir(t)
	defer os.RemoveAll(remote)

	_, err := exec.Command("git", "init", "--bare", remote).CombinedOutput()
	assert.NoError(t, err)

	// repo 1:

	g1, dir1 := newTestRepo(t, "main")
	defer os.RemoveAll(dir1)

	err = os.Chdir(dir1)
	assert.NoError(t, err)

	err = g1.SetRemote("origin", remote)
	assert.NoError(t, err)

	file := "file.txt"
	err = os.WriteFile(filepath.Join(dir1, file), []byte{}, 0644)
	assert.NoError(t, err)

	err = g1.CommitFiles([]string{file}, "Initial commit")
	assert.NoError(t, err)

	err = g1.Push("origin", "main")
	assert.NoError(t, err)

	// repo 2:

	g2, dir2 := newTestRepo(t, "")
	defer os.RemoveAll(dir2)

	err = os.Chdir(dir2)
	assert.NoError(t, err)

	err = g2.SetRemote("origin", remote)
	assert.NoError(t, err)

	err = g2.Pull("origin", "main")
	assert.NoError(t, err)

	// repo1:
	err = os.Chdir(dir1)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(dir1, file), []byte("test\ntest\ntest"), 0644)
	assert.NoError(t, err)

	err = g1.CommitFiles([]string{file}, "Some change")
	assert.NoError(t, err)

	err = g1.Push("origin", "main")
	assert.NoError(t, err)

	// repo2:
	err = os.Chdir(dir2)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(dir2, file), []byte("test\ntext\ntest"), 0644)
	assert.NoError(t, err)

	err = g2.MergeStash("origin", "main", "Merged changes")
	assert.NoError(t, err)
}
func TestNewTag(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	file := "test-file.txt"

	err := os.WriteFile(filepath.Join(dir, file), []byte("test"), 0644)
	assert.NoError(t, err)

	err = g.CommitFiles([]string{file}, "Initial commit")
	assert.NoError(t, err)

	tag := "v0.1.0"
	err = g.NewTag(tag, "Initial release")
	assert.NoError(t, err)

	actual, err := g.LatestTag(true)
	assert.NoError(t, err)

	assert.Equal(t, tag, actual)
}

func TestMoveToRootDir(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	err := os.Mkdir("test", 0755)
	assert.NoError(t, err)

	subdir := filepath.Join(dir, "test")
	err = os.Chdir(subdir)
	assert.NoError(t, err)

	restore := g.MustMoveToRootDir()

	tl, err := os.Getwd()
	assert.NoError(t, err)

	assert.Equal(t, dir, tl)

	err = restore()
	assert.NoError(t, err)

	actual, err := os.Getwd()
	assert.Equal(t, subdir, actual)
}

func TestPull(t *testing.T) {
	remote := tests.NewTempDir(t)
	defer os.RemoveAll(remote)

	_, err := exec.Command("git", "init", "--bare", remote).CombinedOutput()
	assert.NoError(t, err)

	g1, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	err = g1.SetRemote("origin", remote)
	assert.NoError(t, err)

	file := "file.txt"
	err = os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
	assert.NoError(t, err)

	err = g1.CommitFiles([]string{file}, "Initial commit")
	assert.NoError(t, err)

	err = g1.Push("origin", "main")
	assert.NoError(t, err)

	hash, err := g1.LatestHash(true)
	assert.NoError(t, err)

	g2, dir2 := newTestRepo(t, "")
	defer os.RemoveAll(dir2)

	err = g2.SetRemote("origin", remote)
	assert.NoError(t, err)

	err = g2.Pull("origin", "main")
	assert.NoError(t, err)

	actual, err := g2.LatestHash(true)
	assert.NoError(t, err)

	assert.Equal(t, hash, actual)
}

func TestPush(t *testing.T) {
	remote := tests.NewTempDir(t)
	defer os.RemoveAll(remote)

	_, err := exec.Command("git", "init", "--bare", remote).CombinedOutput()
	assert.NoError(t, err)

	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	err = g.SetRemote("origin", remote)
	assert.NoError(t, err)

	remotes, err := g.Remotes()
	assert.NoError(t, err)

	for k, v := range remotes {
		if k != "origin" && v != "git@localhost" {
			t.Fatalf("expected `origin:%s` but got `%s:%s`", remote, k, v)
		}
	}

	file := "file.txt"
	err = os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
	assert.NoError(t, err)

	err = g.CommitFiles([]string{file}, "Initial commit")
	assert.NoError(t, err)

	err = g.Push("origin", "main")
	assert.NoError(t, err)

	err = g.SetUpstreamBranchTo("origin", "main")
	assert.NoError(t, err)

	v, err := g.Config("branch.main.remote")
	assert.NoError(t, err)

	assert.Equal(t, "origin", v)
}

func TestRemoveFromStaging(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	file := "test-file.txt"

	err := os.WriteFile(filepath.Join(dir, file), []byte("test"), 0644)
	assert.NoError(t, err)

	err = g.AddToStaging([]string{file})
	assert.NoError(t, err)

	staged, _, _, err := g.Status()
	assert.NoError(t, err)

	assert.Equal(t, []string{file}, staged)

	err = g.RemoveFromStaging([]string{file}, false)
	assert.NoError(t, err)

	staged, _, _, err = g.Status()
	assert.NoError(t, err)

	assert.Equal(t, 0, len(staged))
}

func TestSetConfig(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	remote := "origin"
	err := g.SetConfig("branch.main.remote", remote)
	assert.NoError(t, err)

	actual, err := g.Config("branch.main.remote")
	assert.NoError(t, err)

	assert.Equal(t, remote, actual)
}

func TestSetRemote(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	err := g.SetRemote("origin", "git@localhost")
	assert.NoError(t, err)

	remotes, err := g.Remotes()
	assert.NoError(t, err)

	for k, v := range remotes {
		if k != "origin" && v != "git@localhost" {
			t.Fatalf("expected `origin:git@localhost` but got `%s:%s`", k, v)
		}
	}
}

func TestStash(t *testing.T) {
	g, dir := newTestRepo(t, "main")
	defer os.RemoveAll(dir)

	file1 := "file1.txt"
	err := os.WriteFile(filepath.Join(dir, file1), []byte{}, 0644)
	assert.NoError(t, err)

	err = g.CommitFiles([]string{file1}, "Initial commit")
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(dir, file1), []byte("test"), 0644)
	assert.NoError(t, err)

	file2 := "file2.txt"
	err = os.WriteFile(filepath.Join(dir, file2), []byte{}, 0644)
	assert.NoError(t, err)

	_, unstaged, untracked, err := g.Status()
	assert.NoError(t, err)

	assert.Equal(t, []string{file1}, unstaged)
	assert.Equal(t, []string{file2}, untracked)

	_, err = g.Stash("", true)
	assert.NoError(t, err)

	_, unstaged, untracked, err = g.Status()
	assert.NoError(t, err)

	assert.Equal(t, 0, len(unstaged))
	assert.Equal(t, 0, len(untracked))
}

func newTestRepo(t *testing.T, initialBranch string) (g Handler, dir string) {
	dir = tests.NewTempDir(t)

	err := os.Chdir(dir)
	assert.NoError(t, err)

	g, err = NewHandler(dir)
	assert.NoError(t, err)

	err = g.Init(initialBranch)
	assert.NoError(t, err)

	return
}
