package logger

import (
	"log"
	"os"
)

type logLevel int

const (
	error logLevel = iota
	info	
	debug	 
)

type Logger struct {
	level logLevel
	logger *log.Logger
}

func New(level string) *Logger {
	var logger = Logger{}
	switch level {
	case "DEBUG" :
		logger.level = debug
	case "INFO" :
		logger.level = info
	case "ERROR" :
		logger.level = error
	default:
		logger.level = info
	}
	logger.logger = log.New(os.Stdout, "", log.LstdFlags)
	return &logger
}

func (l *Logger) Info(msg string) {
	if l.level >= 1 {
		l.logger.SetPrefix("[INFO]")
		l.logger.Println(msg)
	}
}

func (l *Logger) Error(msg string) {
	if l.level >= 0 {
		l.logger.SetPrefix("[ERROR]")
		l.logger.Println(msg)
	}
}

func (l *Logger) Debug(msg string) {
	if l.level >= 2 {
		l.logger.SetPrefix("[DEBUG]")
		l.logger.Println(msg)
	}
}
