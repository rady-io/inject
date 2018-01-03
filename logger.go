package rhapsody

import (
	"github.com/op/go-logging"
	"os"
)

// Logger is a logger base on `github.com/op/go-logging`
type Logger struct {
	*logging.Logger
}

// NewLogger is the factory function of Logger
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
