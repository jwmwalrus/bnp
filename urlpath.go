package bnp

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// PathToURL converts a path to a URL string
func PathToURL(path string) (s string, err error) {
	var u *url.URL
	if u, err = url.Parse(filepath.Clean(strings.ReplaceAll(path, "%", "%25"))); err != nil {
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
	path, u, err := decomposeURL(s)
	if err != nil {
		return
	}

	if u.Scheme != "file" {
		exists = true
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
	path, _, err = decomposeURL(s)
	return
}

func decomposeURL(s string) (path string, u *url.URL, err error) {
	if u, err = url.Parse(s); err != nil {
		return
	}

	if path, err = url.PathUnescape(u.Path); err != nil {
		return
	}

	if u.Scheme == "" {
		u.Scheme = "file"
	}

	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	if u.Fragment != "" {
		path += "#" + u.Fragment
	}

	return
}
