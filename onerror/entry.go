package onerror

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/jwmwalrus/bnp/httpstatus"
)

type Recorder interface {
	LogHTTP(err error, r *http.Response, doNotCloseBody bool) error
	Log(error)
	Panic(error)
	Warn(error)
	With(...slog.Attr) Recorder
}

func NewRecorder(l *slog.Logger) Recorder {
	return &entry{l}
}

type entry struct {
	*slog.Logger
}

// Fatal asserts that no error was given
func (e *entry) Fatal(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		e.withCaller(file, line).(*entry).
			Logger.Error(err.Error())
		os.Exit(1)
	}
}

// LogHTTP logs and HTTP-related error
func (e *entry) LogHTTP(err error, r *http.Response, doNotCloseBody bool) error {
	if err != nil {
		if r != nil {
			e.withStatus(r.StatusCode, r.Status).(*entry).
				Logger.Error(err.Error())
		} else {
			e.Logger.Error(err.Error())
		}
		return err
	} else if r != nil && httpstatus.IsError(r) {
		if !doNotCloseBody {
			defer r.Body.Close()
		}

		var b []byte
		b, _ = io.ReadAll(r.Body)
		msg := string(b)

		e.withStatus(r.StatusCode, r.Status, msg).(*entry).
			Logger.Error(err.Error())
		return fmt.Errorf("ERROR: %v\n\t%v", r.Status, msg)
	}
	return nil
}

// Log logs an error
func (e *entry) Log(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		e.withCaller(file, line).(*entry).
			Logger.Error(err.Error())
	}
}

// Panic asserts that no error was given
func (e *entry) Panic(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		e.withCaller(file, line).(*entry).
			Logger.Error(err.Error())
		panic(err)
	}
}

// Warn warns on error
func (e *entry) Warn(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		e.withCaller(file, line).(*entry).
			Logger.Warn(err.Error())
	}
}

// WithFields adds fields to the entry
func (e *entry) With(a ...slog.Attr) Recorder {
	return &entry{e.Logger.With(a)}
}

func (e *entry) withCaller(file string, line int) Recorder {
	a := callerAttrs(file, line)
	return e.With(a...)
}

func (e *entry) withStatus(statusCode int, status string, msg ...string) Recorder {
	a := statusAttrs(statusCode, status, msg...)
	return e.With(a...)
}

func callerAttrs(file string, line int) []slog.Attr {
	return []slog.Attr{
		slog.String("caller", file),
		slog.Int("callerLine", line),
	}
}

func statusAttrs(statusCode int, status string, msg ...string) []slog.Attr {
	attr := []slog.Attr{
		slog.Int("statusCode", statusCode),
		slog.String("status", status),
	}
	if len(msg) == 1 {
		attr = append(attr, slog.String("error", msg[0]))
	} else if len(msg) > 0 {
		for i, m := range msg {
			attr = append(attr, slog.String("msg"+strconv.Itoa(i+1), m))
		}
	}
	return attr
}

func withCaller(file string, line int) *entry {
	a := callerAttrs(file, line)
	return &entry{slog.With(a)}
}

func withStatus(statusCode int, status string, msg ...string) *entry {
	a := statusAttrs(statusCode, status, msg...)
	return &entry{slog.With(a)}
}
