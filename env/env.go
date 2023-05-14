package env

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// SetDirs Ensure that environment directories exist
func SetDirs(cacheDir, configDir, dataDir, runtimeDir string) (err error) {
	return CreateTheseDirs([]string{
		cacheDir,
		configDir,
		dataDir,
		runtimeDir,
	})
}

// CreateTheseDirs creates the non-existent directories from the given list
//
// This function will return at the first failure to create any of the
// given directories.
func CreateTheseDirs(dirs []string) (err error) {
	for i, d := range dirs {
		if _, err = os.Stat(d); os.IsNotExist(err) {
			if err = os.MkdirAll(d, 0755); err != nil {
				err = fmt.Errorf("failed to create dirs[%d] = %s: %w", i, d, err)
				return
			}
		}
	}

	return
}

// FindExec finds the full path for the given binary
func FindExec(bin string) string {
	if path := findBinary(bin); path != "" {
		return path
	}

	if full, err := exec.LookPath(bin); !os.IsNotExist(err) {
		return full
	}

	return ""
}

// FindLibExec finds the full path for the given binary in libexec directories
func FindLibExec(bin, app string) string {
	if path := findBinary(filepath.Join(app, bin)); path != "" {
		return path
	}

	list := []string{"/usr/local/libexec", "/usr/libexec", "/libexec"}

	pathExists := func(path string) bool {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return true
		}
		return false
	}

	for _, l := range list {
		path := filepath.Join(l, app, bin)
		if pathExists(path) {
			path, _ = filepath.Abs(path)
			return path
		}
	}
	return ""
}

func findBinary(bin string) string {
	pathExists := func(path string) bool {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return true
		}
		return false
	}

	path := filepath.Join(".", bin)
	for {
		if pathExists(path) {
			path, _ = filepath.Abs(path)
			return path
		}
		if apath, _ := filepath.Abs(path); apath == filepath.Join(string(filepath.Separator), bin) {
			break
		}
		path = filepath.Join("..", path)
	}

	return ""
}
