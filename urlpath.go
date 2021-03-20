package bnp

import (
	"net/url"
	"os"
	"path/filepath"
)

// PathToURL converts a path to a URL string
func PathToURL(path string) (s string, err error) {
	var u *url.URL
	if u, err = url.Parse(filepath.Clean(path)); err != nil {
		return
	}

	if u.Scheme == "" {
		u.Scheme = "file"
	}

	s = u.String()
	return
}

// URLExists checks if the string corresponds to an existing location
// It always returns true if the Scheme is not file
func URLExists(s string) (exists bool) {
	var u *url.URL
	var err error
	if u, err = url.Parse(s); err != nil {
		return
	}

	if u.Scheme != "file" {
		exists = true
		return
	}

	var path string
	if path, err = url.PathUnescape(u.Path); err != nil {
		return
	}
	if _, err = os.Stat(path); !os.IsNotExist(err) {
		exists = true
		return
	}

	return
}

// URLToPath convertes the given URL string to a file path
func URLToPath(s string) (path string, err error) {
	var u *url.URL
	if u, err = url.Parse(s); err != nil {
		return
	}

	path, err = url.PathUnescape(u.Path)
	return
}
