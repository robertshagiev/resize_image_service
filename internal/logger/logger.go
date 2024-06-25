package logger

import (
	"log"
	"os"
)

type logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}

func New() Logger {
	return &logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *logger) Info(msg string) {
	l.infoLogger.Println(msg)
}

func (l *logger) Error(msg string) {
	l.errorLogger.Println(msg)
}

func (l *logger) Fatal(msg string) {
	l.errorLogger.Fatal(msg)
}
