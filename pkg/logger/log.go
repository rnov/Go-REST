package logger

import (
	"fmt"
	"github.com/rnov/Go-REST/pkg/config"
	"os"

	"github.com/op/go-logging"
)

// note: get from go-logging/log.go
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

// done
func NewLogger(log config.LoggerConfig, path string) *logging.Logger {

	// initialice a logger
	logger := logging.MustGetLogger(log.Name)

	loggerPath := fmt.Sprintf("%s%s", path, log.File)
	logFile, err := os.OpenFile(loggerPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{errors}`,
	)

	backend := logging.NewLogBackend(logFile, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(backendLeveled, backendFormatter)

	return logger
}
