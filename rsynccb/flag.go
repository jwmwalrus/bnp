package rsynccb

import (
	"fmt"
	"log/slog"
	"strings"
)

type rsyncFlagType int

const (
	rsyncBoolFlag rsyncFlagType = iota
	rsyncStringFlag
)

type rsyncFlagValue struct {
	ftype rsyncFlagType
	value any
}

type rsyncFlags map[string]rsyncFlagValue

var (
	supportedFlags rsyncFlags
)

func init() {
	supportedFlags = rsyncFlags{
		"verbose":             rsyncFlagValue{ftype: rsyncBoolFlag},
		"info":                rsyncFlagValue{ftype: rsyncStringFlag},
		"debug":               rsyncFlagValue{ftype: rsyncStringFlag},
		"stderr":              rsyncFlagValue{ftype: rsyncStringFlag},
		"quiet":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"no-motd":             rsyncFlagValue{ftype: rsyncBoolFlag},
		"checksum":            rsyncFlagValue{ftype: rsyncBoolFlag},
		"archive":             rsyncFlagValue{ftype: rsyncBoolFlag},
		"recursive":           rsyncFlagValue{ftype: rsyncBoolFlag},
		"relative":            rsyncFlagValue{ftype: rsyncBoolFlag},
		"no-implied-dirs":     rsyncFlagValue{ftype: rsyncBoolFlag},
		"backup":              rsyncFlagValue{ftype: rsyncBoolFlag},
		"backup-dir":          rsyncFlagValue{ftype: rsyncStringFlag},
		"suffix":              rsyncFlagValue{ftype: rsyncStringFlag},
		"update":              rsyncFlagValue{ftype: rsyncBoolFlag},
		"inplace":             rsyncFlagValue{ftype: rsyncBoolFlag},
		"append":              rsyncFlagValue{ftype: rsyncBoolFlag},
		"append-verify":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"dirs":                rsyncFlagValue{ftype: rsyncBoolFlag},
		"old-dirs":            rsyncFlagValue{ftype: rsyncBoolFlag},
		"old-d":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"mkpath":              rsyncFlagValue{ftype: rsyncBoolFlag},
		"links":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"copy-links":          rsyncFlagValue{ftype: rsyncBoolFlag},
		"copy-unsafe-links":   rsyncFlagValue{ftype: rsyncBoolFlag},
		"safe-links":          rsyncFlagValue{ftype: rsyncBoolFlag},
		"munge-links":         rsyncFlagValue{ftype: rsyncBoolFlag},
		"copy-dirlinks":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"keep-dirlinks":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"hard-links":          rsyncFlagValue{ftype: rsyncBoolFlag},
		"perms":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"executability":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"chmod":               rsyncFlagValue{ftype: rsyncStringFlag},
		"acls":                rsyncFlagValue{ftype: rsyncBoolFlag},
		"xattrs":              rsyncFlagValue{ftype: rsyncBoolFlag},
		"owner":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"group":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"devices":             rsyncFlagValue{ftype: rsyncBoolFlag},
		"copy-devices":        rsyncFlagValue{ftype: rsyncBoolFlag},
		"write-devices":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"specials":            rsyncFlagValue{ftype: rsyncBoolFlag},
		"times":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"atimes":              rsyncFlagValue{ftype: rsyncBoolFlag},
		"open-noatime":        rsyncFlagValue{ftype: rsyncBoolFlag},
		"crtimes":             rsyncFlagValue{ftype: rsyncBoolFlag},
		"omit-dir-times":      rsyncFlagValue{ftype: rsyncBoolFlag},
		"omit-link-times":     rsyncFlagValue{ftype: rsyncBoolFlag},
		"super":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"fake-super":          rsyncFlagValue{ftype: rsyncBoolFlag},
		"sparse":              rsyncFlagValue{ftype: rsyncBoolFlag},
		"preallocate":         rsyncFlagValue{ftype: rsyncBoolFlag},
		"dry-run":             rsyncFlagValue{ftype: rsyncBoolFlag},
		"whole-file":          rsyncFlagValue{ftype: rsyncBoolFlag},
		"checksum-choice":     rsyncFlagValue{ftype: rsyncStringFlag},
		"one-file-system":     rsyncFlagValue{ftype: rsyncBoolFlag},
		"block-size":          rsyncFlagValue{ftype: rsyncStringFlag},
		"rsh":                 rsyncFlagValue{ftype: rsyncStringFlag},
		"rsync-path":          rsyncFlagValue{ftype: rsyncStringFlag},
		"existing":            rsyncFlagValue{ftype: rsyncBoolFlag},
		"ignore-existing":     rsyncFlagValue{ftype: rsyncBoolFlag},
		"remove-source-files": rsyncFlagValue{ftype: rsyncBoolFlag},
		"del":                 rsyncFlagValue{ftype: rsyncBoolFlag},
		"delete":              rsyncFlagValue{ftype: rsyncBoolFlag},
		"delete-before":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"delete-during":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"delete-delay":        rsyncFlagValue{ftype: rsyncBoolFlag},
		"delete-after":        rsyncFlagValue{ftype: rsyncBoolFlag},
		"delete-excluded":     rsyncFlagValue{ftype: rsyncBoolFlag},
		"ignore-missing-args": rsyncFlagValue{ftype: rsyncBoolFlag},
		"delete-missing-args": rsyncFlagValue{ftype: rsyncBoolFlag},
		"ignore-errors":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"force":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"max-delete":          rsyncFlagValue{ftype: rsyncStringFlag},
		"max-size":            rsyncFlagValue{ftype: rsyncStringFlag},
		"min-size":            rsyncFlagValue{ftype: rsyncStringFlag},
		"max-alloc":           rsyncFlagValue{ftype: rsyncStringFlag},
		"partial":             rsyncFlagValue{ftype: rsyncBoolFlag},
		"partial-dir":         rsyncFlagValue{ftype: rsyncStringFlag},
		"delay-updates":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"prune-empty-dirs":    rsyncFlagValue{ftype: rsyncBoolFlag},
		"numeric-ids":         rsyncFlagValue{ftype: rsyncBoolFlag},
		"usermap":             rsyncFlagValue{ftype: rsyncStringFlag},
		"groupmap":            rsyncFlagValue{ftype: rsyncStringFlag},
		"chown":               rsyncFlagValue{ftype: rsyncStringFlag},
		"timeout":             rsyncFlagValue{ftype: rsyncStringFlag},
		"contimeout":          rsyncFlagValue{ftype: rsyncStringFlag},
		"ignore-times":        rsyncFlagValue{ftype: rsyncBoolFlag},
		"size-only":           rsyncFlagValue{ftype: rsyncBoolFlag},
		"modify-window":       rsyncFlagValue{ftype: rsyncStringFlag},
		"temp-dir":            rsyncFlagValue{ftype: rsyncStringFlag},
		"fuzzy":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"compare-dest":        rsyncFlagValue{ftype: rsyncStringFlag},
		"copy-dest":           rsyncFlagValue{ftype: rsyncStringFlag},
		"link-dest":           rsyncFlagValue{ftype: rsyncStringFlag},
		"compress":            rsyncFlagValue{ftype: rsyncBoolFlag},
		"compress-choice":     rsyncFlagValue{ftype: rsyncStringFlag},
		"compress-level":      rsyncFlagValue{ftype: rsyncStringFlag},
		"skip-compress":       rsyncFlagValue{ftype: rsyncStringFlag},
		"cvs-exclude":         rsyncFlagValue{ftype: rsyncBoolFlag},
		"filter":              rsyncFlagValue{ftype: rsyncStringFlag},
		"exclude":             rsyncFlagValue{ftype: rsyncStringFlag},
		"exclude-from":        rsyncFlagValue{ftype: rsyncStringFlag},
		"include":             rsyncFlagValue{ftype: rsyncStringFlag},
		"include-from":        rsyncFlagValue{ftype: rsyncStringFlag},
		"files-from":          rsyncFlagValue{ftype: rsyncStringFlag},
		"from0":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"old-args":            rsyncFlagValue{ftype: rsyncBoolFlag},
		"secluded-args":       rsyncFlagValue{ftype: rsyncBoolFlag},
		"trust-sender":        rsyncFlagValue{ftype: rsyncBoolFlag},
		"copy-as":             rsyncFlagValue{ftype: rsyncStringFlag},
		"address":             rsyncFlagValue{ftype: rsyncStringFlag},
		"port":                rsyncFlagValue{ftype: rsyncStringFlag},
		"sockopts":            rsyncFlagValue{ftype: rsyncStringFlag},
		"blocking-io":         rsyncFlagValue{ftype: rsyncBoolFlag},
		"outbuf":              rsyncFlagValue{ftype: rsyncStringFlag},
		"stats":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"8-bit-output":        rsyncFlagValue{ftype: rsyncBoolFlag},
		"human-readable":      rsyncFlagValue{ftype: rsyncBoolFlag},
		"progress":            rsyncFlagValue{ftype: rsyncBoolFlag},
		"itemize-changes":     rsyncFlagValue{ftype: rsyncBoolFlag},
		"remote-option":       rsyncFlagValue{ftype: rsyncStringFlag},
		"out-format":          rsyncFlagValue{ftype: rsyncStringFlag},
		"log-file":            rsyncFlagValue{ftype: rsyncStringFlag},
		"log-file-format":     rsyncFlagValue{ftype: rsyncStringFlag},
		"password-file":       rsyncFlagValue{ftype: rsyncStringFlag},
		"early-input":         rsyncFlagValue{ftype: rsyncStringFlag},
		"list-only":           rsyncFlagValue{ftype: rsyncBoolFlag},
		"bwlimit":             rsyncFlagValue{ftype: rsyncStringFlag},
		"stop-after":          rsyncFlagValue{ftype: rsyncStringFlag},
		"stop-at":             rsyncFlagValue{ftype: rsyncStringFlag},
		"fsync":               rsyncFlagValue{ftype: rsyncBoolFlag},
		"write-batch":         rsyncFlagValue{ftype: rsyncStringFlag},
		"only-write-batch":    rsyncFlagValue{ftype: rsyncStringFlag},
		"read-batch":          rsyncFlagValue{ftype: rsyncStringFlag},
		"protocol":            rsyncFlagValue{ftype: rsyncStringFlag},
		"iconv":               rsyncFlagValue{ftype: rsyncStringFlag},
		"checksum-seed":       rsyncFlagValue{ftype: rsyncStringFlag},
		"ipv4":                rsyncFlagValue{ftype: rsyncBoolFlag},
		"ipv6":                rsyncFlagValue{ftype: rsyncBoolFlag},
		"version":             rsyncFlagValue{ftype: rsyncBoolFlag},
		"daemon":              rsyncFlagValue{ftype: rsyncBoolFlag},
		"config":              rsyncFlagValue{ftype: rsyncStringFlag},
		"dparam":              rsyncFlagValue{ftype: rsyncStringFlag},
		"no-detach":           rsyncFlagValue{ftype: rsyncBoolFlag},
	}
}

func (f rsyncFlags) set(name string, value any) error {
	v, ok := f[name]
	if !ok {
		return fmt.Errorf("unsupported flag: %v", name)
	}

	switch v.ftype {
	case rsyncBoolFlag:
		b, ok := value.(bool)
		if !ok {
			return fmt.Errorf("boolean value expected for flag: %v", name)
		}
		v.value = b
	case rsyncStringFlag:
		s, ok := value.(string)
		if !ok {
			s, ok := value.([]string)
			if !ok {
				return fmt.Errorf("string or []string value expected for flag: %v", name)
			}
			v.value = s
		} else {
			v.value = []string{s}
		}
	}

	f[name] = v

	return nil
}

func (f rsyncFlags) unset(name string) {
	v, ok := f[name]
	if !ok {
		slog.Error("Unsupported flag", "flag", name)
		return
	}

	v.value = nil
	f[name] = v
}

func (f rsyncFlags) isSet(name string) bool {
	v, ok := f[name]
	if !ok {
		slog.Error("Unsupported flag", "flag", name)
		return false
	}

	if v.value == nil {
		return false
	}

	switch v.ftype {
	case rsyncBoolFlag:
		b, ok := v.value.(bool)
		if !ok {
			slog.With(
				"flag", name,
				"value", v.value,
			).Error("Unexpected value for flag")
			return false
		}
		return b
	case rsyncStringFlag:
		s, ok := v.value.([]string)
		if !ok {
			slog.With(
				"flag", name,
				"value", v.value,
			).Error("Unexpected value for flag")
			return false
		}
		return len(s) > 0
	}

	return false
}

func (f rsyncFlags) Format(name string) string {
	v, ok := f[name]
	if !ok {
		slog.Error("Unsupported flag", "flag", name)
		return ""
	}

	if v.value == nil {
		return ""
	}

	switch v.ftype {
	case rsyncBoolFlag:
		b, ok := v.value.(bool)
		if !ok {
			slog.With(
				"flag", name,
				"value", v.value,
			).Error("Unexpected value for flag")
			return ""
		}

		if b {
			return "--" + name
		}
	case rsyncStringFlag:
		s, ok := v.value.([]string)
		if !ok {
			slog.With(
				"flag", name,
				"value", v.value,
			).Error("Unexpected value for flag")
			return ""
		}

		res := ""
		for _, x := range s {
			res += " --" + name + "=" + x
		}
		return strings.TrimSpace(res)
	}

	return ""
}
