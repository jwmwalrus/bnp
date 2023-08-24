package cron

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	cronDir = "/etc/cron.d"
)

type entry struct {
	c    *Cron
	user string
	cmd  string
}

// Install builds and installs a crontab
type Install struct {
	comment string
	entries []entry
}

// NewInstall returns a new crontab installer
func NewInstall(comment string) *Install {
	return &Install{comment: comment}
}

// Add adds entry to the installer
func (i *Install) Add(c *Cron, user, cmd string) error {
	i.entries = append(i.entries, entry{c, user, cmd})
	return nil
}

// AsUser installs the crontab as user.
// It invokes the crontab command for the operation.
func (i *Install) AsUser() error {

	// how to ensure other entries are not affected?
	return nil
}

// AsSystem installs the entries to the system-wide crontab
func (i *Install) AsSystem(name string) error {
	cronLines := ""

	if i.comment != "" {
		cronLines += "#" + i.comment + "\n"
	}

	for _, e := range i.entries {
		l := e.c.Format(e.user, e.cmd)
		cronLines += l + "\n"
	}

	if _, err := os.Stat(cronDir); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("The intended cron file contents are:\n\n%s\n", cronLines)
		return fmt.Errorf("directory `%s` not found", cronDir)
	}

	return os.WriteFile(filepath.Join(cronDir, name), []byte(cronLines), 0644)
}
