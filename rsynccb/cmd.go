package rsynccb

import "os/exec"

// RsyncCmdBuilder defines the rsync command builder
type RsyncCmdBuilder struct {
	flags rsyncFlags
}

// New return a new RsyncCmd instance
func New() RsyncCmdBuilder {
	r := RsyncCmdBuilder{}
	r.flags = supportedFlags
	return r
}

// Build builds the rsync command
func (r *RsyncCmdBuilder) Build(src, dest string) *exec.Cmd {
	cmd := "rsync"

	args := []string{}
	for k := range r.flags {
		if !r.flags.isSet(k) {
			continue
		}
		args = append(args, r.flags.Format(k))
	}
	args = append(args, src, dest)

	return exec.Command(cmd, args...)
}

// Set sets the value of a flag
func (r *RsyncCmdBuilder) Set(name string, value any) error {
	return r.flags.set(name, value)
}

// Unset unsets the value of a flag
func (r *RsyncCmdBuilder) Unset(name string) {
	r.flags.unset(name)
}

// IsSet returns true if the flag is set
func (r *RsyncCmdBuilder) IsSet(name string) bool {
	return r.flags.isSet(name)
}
