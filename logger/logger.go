package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log = logrus.New()

func init() {
	fmt := prefixed.TextFormatter{
		FullTimestamp: true,
	}
	log.Formatter = &fmt
	log.Level = logrus.DebugLevel
}

func Info(message string) {
	log.Info(message)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Debug(message string) {
	log.Debug(message)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Error(message string) {
	log.Error(message)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
