package bnp

import (
	"io"
	"os"

	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ParseArgs parses the given command line arguments
func ParseArgs(logFile *lumberjack.Logger, echoLogging, verbose *bool, severity *string) (args []string) {
	getopt.Parse()
	args = getopt.Args()
	arg0 := []string{os.Args[0]}
	args = append(arg0, args...)

	resolveSeverity(severity, verbose)

	if *echoLogging {
		mw := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(mw)
	}

	return
}

// SetEnvDirs Ensure that environment directories exist
func SetEnvDirs(configFile, lockFile, cacheDir, configDir, dataDir, runtimeDir string) (err error) {

	if _, err = os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(cacheDir, 0755)
		if err != nil {
			return
		}
	}

	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return
		}
	}

	if _, err = os.Stat(dataDir); os.IsNotExist(err) {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			return
		}
	}

	if _, err = os.Stat(runtimeDir); os.IsNotExist(err) {
		err = os.MkdirAll(runtimeDir, 0755)
	}

	return
}

func resolveSeverity(severity *string, verbose *bool) {
	givenSeverity := *severity

	if givenSeverity == "" {
		if *verbose {
			*severity = "info"
		} else {
			*severity = "error"
		}
	} else {
		if _, err := log.ParseLevel(givenSeverity); err != nil {
			*severity = "error"
		} else {
			*severity = givenSeverity
		}
	}

	level, _ := log.ParseLevel(*severity)
	log.SetLevel(level)
	log.SetReportCaller(*severity == "debug")

	return
}
