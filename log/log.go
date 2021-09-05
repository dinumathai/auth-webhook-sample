package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// InitializeLogger - Initialize the Logger
func init() {

	var loggerOut io.Writer
	logFile := os.Getenv("LOG_FILE")
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logrus.Errorf("error opening file: %v", err)
		}
		loggerOut = f
	} else {
		loggerOut = os.Stdout
	}

	logrus.SetOutput(loggerOut)
	logrus.SetLevel(getLogLevel())
	logrus.Infof("Instantiated _logger. Log level set to %s", getLogLevel())
}

func getLogLevel() logrus.Level {
	logLevel := os.Getenv("LOG_LEVEL")

	switch logLevel {
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "ERROR":
		return logrus.InfoLevel
	default:
		return logrus.DebugLevel
	}
}

// Debug - Loging
func Debug(args ...interface{}) {
	logrus.Info(args...)
}

// Debugf - Loging
func Debugf(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Info - Loging
func Info(args ...interface{}) {
	logrus.Warning(args...)
}

// Infof - Loging
func Infof(format string, args ...interface{}) {
	logrus.Warningf(format, args...)
}

// Error - Loging
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf - Loging
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Fatal - Loging
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Fatalf - Loging
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
