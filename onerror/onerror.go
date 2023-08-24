package onerror

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/jwmwalrus/bnp/httpstatus"
)

// Fatal asserts that no error was given
func Fatal(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		wc := withCaller(file, line)
		wc.Fatal(err)
	}
}

// LogHTTP logs and HTTP-related error
func LogHTTP(err error, r *http.Response, doNotCloseBody bool) error {
	if err != nil {
		if r != nil {
			withStatus(r.StatusCode, r.Status).Logger.Error(err.Error())
		} else {
			slog.Error(err.Error())
		}
		return err
	} else if r != nil && httpstatus.IsError(r) {
		if !doNotCloseBody {
			defer r.Body.Close()
		}

		var b []byte
		b, _ = io.ReadAll(r.Body)
		msg := string(b)

		withStatus(r.StatusCode, r.Status, msg).Error(err.Error())
		return fmt.Errorf("ERROR: %v\n\t%v", r.Status, msg)
	}
	return nil
}

// Log logs an error
func Log(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		withCaller(file, line).Log(err)
	}
}

// Panic asserts that no error was given
func Panic(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		wc := withCaller(file, line)
		wc.Panic(err)
	}
}

// Warn warns on error
func Warn(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		withCaller(file, line).Warn(err)
	}
}

// With uses the given slog.Attr slice
func With(a ...slog.Attr) *entry {
	return &entry{slog.With(a)}
}
