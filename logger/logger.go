package logger

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Level int

const (
	DEBUG Level = 0
	INFO  Level = 1
	WARN  Level = 2
	ERROR Level = 3
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

func parseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

type Logger struct {
	module string
	level  Level
	logger *log.Logger
}

func New(w io.Writer, level string) *Logger {
	return &Logger{
		logger: log.New(w, "", log.LstdFlags|log.Lshortfile),
		level:  parseLevel(level),
	}
}

func (l *Logger) NewModule(module string) *Logger {
	return &Logger{
		module: module,
		logger: l.logger,
		level:  l.level,
	}
}

func (l *Logger) log(level Level, format string, v ...any) {
	if level < l.level {
		return
	}
	prefix := fmt.Sprintf("[%s] [%s]: ", l.module, levelNames[level])
	msg := fmt.Sprintf(format, v...)
	l.logger.Output(3, prefix+msg)
}

func (l *Logger) Debug(format string, v ...any) { l.log(DEBUG, format, v...) }
func (l *Logger) Info(format string, v ...any)  { l.log(INFO, format, v...) }
func (l *Logger) Warn(format string, v ...any)  { l.log(WARN, format, v...) }
func (l *Logger) Error(format string, v ...any) { l.log(ERROR, format, v...) }
func (l *Logger) Panic(format string, v ...any) {
	l.log(ERROR, format, v...)
	log.Panic()
}
