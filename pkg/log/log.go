package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)
	log.Formatter = &logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
}

func Info(message ...interface{}) {
	log.Info(message...)
}

func Trace(message ...interface{}) {
	log.Trace(message...)
}

func Debug(message ...interface{}) {
	log.Debug(message...)
}

func Warn(message ...interface{}) {
	log.Warn(message...)
}

func Error(message ...interface{}) {
	log.Error(message...)
}

func Fatal(message ...interface{}) {
	log.Fatal(message...)
}

func Panic(message ...interface{}) {
	log.Panic(message...)
}

func Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func Tracef(format string, v ...interface{}) {
	log.Tracef(format, v...)
}

func Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}
