package bnp

import (
	"runtime"

	log "github.com/sirupsen/logrus"
)

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
