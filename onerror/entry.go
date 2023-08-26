package onerror

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jwmwalrus/bnp/httpstatus"
)

// Recorder defines the interface to record an error
type Recorder interface {
	Fatal(error)
	LogHTTP(err error, r *http.Response, doNotCloseBody bool) error
	Log(error)
	Panic(error)
	Warn(error)
	With(...any) Recorder
}

// NewRecorder returns a new recorder for the given slog.Logger
func NewRecorder(l *slog.Logger) Recorder {
	return &entry{Logger: l}
}

type entry struct {
	*slog.Logger
	caller []slog.Attr
	pc     uintptr
}

// Fatal asserts that no error was given
func (e *entry) Fatal(err error) {
	if err != nil {
		if len(e.caller) == 0 {
			e.caller, e.pc = callerAttrs()
		}
		e.logError(slog.LevelError, err)
		os.Exit(1)
	}
}

// LogHTTP logs and HTTP-related error
func (e *entry) LogHTTP(err error, r *http.Response, doNotCloseBody bool) error {
	var addRecord bool
	if err != nil {
		addRecord = true
	} else if r != nil && httpstatus.IsError(r) {
		addRecord = true

		if !doNotCloseBody {
			defer r.Body.Close()
		}

		var b []byte
		b, _ = io.ReadAll(r.Body)
		msg := string(b)

		err = fmt.Errorf(msg)
	}

	if !addRecord {
		return nil
	}

	if len(e.caller) == 0 {
		e.caller, e.pc = callerAttrs()
	}

	if r != nil {
		a := statusAttrs(r.StatusCode, r.Status)
		e.Logger = e.Logger.With(a...)
	}

	e.logError(slog.LevelError, err)

	return err
}

// Log logs an error
func (e *entry) Log(err error) {
	if err != nil {
		if len(e.caller) == 0 {
			e.caller, e.pc = callerAttrs()
		}
		e.logError(slog.LevelError, err)
	}
}

// Panic asserts that no error was given
func (e *entry) Panic(err error) {
	if err != nil {
		if len(e.caller) == 0 {
			e.caller, e.pc = callerAttrs()
		}
		e.logError(slog.LevelError, err)
		panic(err)
	}
}

// Warn warns on error
func (e *entry) Warn(err error) {
	if err != nil {
		if len(e.caller) == 0 {
			e.caller, e.pc = callerAttrs()
		}
		e.logError(slog.LevelWarn, err)
	}
}

// WithFields adds fields to the entry
func (e *entry) With(a ...any) Recorder {
	if e.Logger != nil {
		return &entry{Logger: e.Logger.With(a...)}
	}
	return &entry{Logger: slog.Default().With(a...)}
}

func (e *entry) logError(level slog.Level, err error) {
	record := slog.NewRecord(time.Now(), level, "ONERROR-"+err.Error(), e.pc)
	record.AddAttrs(e.caller...)
	e.Logger.Handler().Handle(context.Background(), record)
}

func callerAttrs() ([]slog.Attr, uintptr) {
	pc, file, line, _ := runtime.Caller(2)
	return []slog.Attr{
		slog.String("caller", file),
		slog.Int("callerLine", line),
	}, pc
}

func statusAttrs(statusCode int, status string) []any {
	return []any{
		"status_code", statusCode,
		"status", status,
	}
}
