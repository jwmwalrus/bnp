package onerror

import (
	"log/slog"
	"net/http"
)

// Fatal asserts that no error was given
func Fatal(err error) {
	caller, pc := callerAttrs()
	e := &entry{Logger: slog.Default(), caller: caller, pc: pc}
	e.Fatal(err)
}

// LogHTTP logs an HTTP-related error
func LogHTTP(err error, r *http.Response, doNotCloseBody bool) error {
	caller, pc := callerAttrs()
	e := &entry{Logger: slog.Default(), caller: caller, pc: pc}
	return e.LogHTTP(err, r, doNotCloseBody)
}

// Log logs an error
func Log(err error) {
	caller, pc := callerAttrs()
	e := &entry{Logger: slog.Default(), caller: caller, pc: pc}
	e.Log(err)
}

// Panic asserts that no error was given
func Panic(err error) {
	caller, pc := callerAttrs()
	e := &entry{Logger: slog.Default(), caller: caller, pc: pc}
	e.Panic(err)
}

// Warn warns on error
func Warn(err error) {
	caller, pc := callerAttrs()
	e := &entry{Logger: slog.Default(), caller: caller, pc: pc}
	e.Warn(err)
}

// With uses the given slog.Attr slice
func With(a ...any) Recorder {
	return &entry{Logger: slog.With(a...)}
}
