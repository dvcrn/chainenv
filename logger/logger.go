package logger

import (
	"log"
	"os"
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	err     *log.Logger
	isDebug bool
}

func NewLogger(debugEnabled bool) *Logger {
	return &Logger{
		debug:   log.New(os.Stdout, "DEBUG: ", 0),
		info:    log.New(os.Stdout, "INFO: ", 0),
		err:     log.New(os.Stderr, "ERR: ", 0),
		isDebug: debugEnabled,
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.isDebug {
		l.debug.Printf(format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

func (l *Logger) Err(format string, v ...interface{}) {
	l.err.Printf(format, v...)
}
