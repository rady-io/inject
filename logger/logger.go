package logger

import (
	"github.com/op/go-logging"
	"os"
)

type Logger struct {
	*logging.Logger
}

func NewLogger() *Logger {
	log := logging.MustGetLogger("example")
	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	return &Logger{log}
}
