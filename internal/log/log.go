package log

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

const LogLevelFlagName = "loglevel"

func init() {
	flag.String(
		LogLevelFlagName,
		logrus.InfoLevel.String(),
		fmt.Sprintf("Log level available: %s", logrus.AllLevels),
	)

	loglevel, err := logrus.ParseLevel(flag.Lookup(LogLevelFlagName).Value.String())
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(loglevel)
	}

	if term.IsTerminal(int(os.Stdout.Fd())) {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	logrus.Fatalln(args...)
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	logrus.Panicln(args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

func Print(args ...interface{}) {
	logrus.Print(args...)
}

func Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)
}

func Println(args ...interface{}) {
	logrus.Println(args...)
}

func Trace(args ...interface{}) {
	logrus.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

func Traceln(args ...interface{}) {
	logrus.Traceln(args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

func WithContext(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx)
}

func WithTime(t time.Time) *logrus.Entry {
	return logrus.WithTime(t)
}
