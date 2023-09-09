package onerror

import (
	"log/slog"
	"net/http"
)

// Fatal asserts that no error was given
func Fatal(err error) {
	e := &entry{Logger: slog.Default(), pc: callerPC(2)}
	e.Fatal(err)
}

// LogHTTP logs an HTTP-related error
func LogHTTP(err error, r *http.Response, doNotCloseBody bool) error {
	e := &entry{Logger: slog.Default(), pc: callerPC(2)}
	return e.LogHTTP(err, r, doNotCloseBody)
}

// Log logs an error
func Log(err error) {
	e := &entry{Logger: slog.Default(), pc: callerPC(2)}
	e.Log(err)
}

// Panic asserts that no error was given
func Panic(err error) {
	e := &entry{Logger: slog.Default(), pc: callerPC(2)}
	e.Panic(err)
}

// Warn warns on error
func Warn(err error) {
	e := &entry{Logger: slog.Default(), pc: callerPC(2)}
	e.Warn(err)
}

// With uses the given slog.Attr slice
func With(a ...any) Recorder {
	return &entry{Logger: slog.With(a...)}
}
