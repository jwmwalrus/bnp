package onerror

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/jwmwalrus/bnp/httpstatus"
	log "github.com/sirupsen/logrus"
)

// LogHTTP logs and HTTP-related error
func LogHTTP(err error, r *http.Response, doNotCloseBody bool) error {
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
	} else if r != nil && httpstatus.IsError(r) {
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

// Log logs an error
func Log(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{
			"caller":     file,
			"callerLine": line,
		}).Error(err)
	}
}

// Panic asserts that no error was given
func Panic(err error) {
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

// Warn warns on error
func Warn(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{
			"caller":     file,
			"callerLine": line,
		}).Warn(err)
	}
}
