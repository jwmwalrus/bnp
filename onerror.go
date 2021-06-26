package bnp

import (
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
	"gopkg.in/errgo.v2/errors"
)

// LogHTTPError logs and HTTP-related error
func LogHTTPError(r *http.Response, err error) error {
	if err != nil {
		log.WithFields(log.Fields{
			"statusCode":    r.StatusCode,
			"statusMessage": r.Status,
		}).Error(err)
		return err
	} else {
		log.WithFields(log.Fields{
			"statusCode": r.StatusCode,
		}).Error(r.Status)
		return errors.New(r.Status)
	}
}

// LogOnError logs an error
func LogOnError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{
			"caller":     file,
			"callerLine": line,
		}).Error(err)
	}
}

// PanicOnError asserts that no error was given
func PanicOnError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{
			"caller":     file,
			"callerLine": line,
		}).Error(err)
		log.WithFields(log.Fields{
			"caller":     file,
			"callerLine": line,
		}).Fatal(err)
	}
}

// WarnOnError warns on error
func WarnOnError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{
			"caller":     file,
			"callerLine": line,
		}).Warn(err)
	}
}
