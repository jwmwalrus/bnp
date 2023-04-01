package env

import (
	"fmt"
	"os"
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
