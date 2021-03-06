package logger

import (
	"os"

	"github.com/op/go-logging"
)

type Loggers interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

func NewLogger() *logging.Logger {
	// initialize a logger
	logger := logging.MustGetLogger("goREST")

	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	// set error backend
	backendError := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backendError, format)

	backendLeveledError := logging.AddModuleLevel(backendError)
	backendLeveledError.SetLevel(logging.ERROR, "")

	// set info backend
	backendInfo := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter = logging.NewBackendFormatter(backendInfo, format)

	backendLeveledInfo := logging.AddModuleLevel(backendInfo)
	backendLeveledInfo.SetLevel(logging.INFO, "")

	logging.SetBackend(backendLeveledInfo, backendFormatter)

	return logger
}
