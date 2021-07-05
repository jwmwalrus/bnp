package bnp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// LogHTTPError logs and HTTP-related error
func LogHTTPError(err error, r *http.Response, doNotCloseBody bool) error {
	if err != nil {
		if r != nil {
			log.WithFields(log.Fields{
				"statusCode": r.StatusCode,
				"status":     r.Status,
			}).Error(err)
		} else {
			log.Error(err)
		}
		return err
	} else if r != nil && IsHTTPError(r) {
		if !doNotCloseBody {
			defer r.Body.Close()
		}

		var b []byte
		b, _ = ioutil.ReadAll(r.Body)
		msg := string(b)

		log.WithFields(log.Fields{
			"statusCode": r.StatusCode,
			"status":     r.Status,
			"error":      msg,
		}).Error(r.Status)
		return fmt.Errorf("ERROR: %v\n\t%v", r.Status, msg)
	}
	return nil
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
