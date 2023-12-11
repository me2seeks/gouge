package cio

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	prefix      string
	Info, Debug bool
	logger      *log.Logger
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
		logger: log.New(os.Stderr, "", log.Ldate|log.Ltime),
		Info:   false,
		Debug:  false,
	}
}

func (l *Logger) Fork(prefix string, args ...interface{}) *Logger {
	//slip the parent prefix at the front
	args = append([]interface{}{l.prefix}, args...)
	ll := NewLogger(fmt.Sprintf("%s: "+prefix, args...))
	ll.Info = l.Info
	ll.Debug = l.Debug
	return ll
}

func (l *Logger) Infof(f string, args ...interface{}) {
	if l.Info {
		l.logger.Printf(l.prefix+": "+f, args...)
	}
}

func (l *Logger) Debugf(f string, args ...interface{}) {
	if l.Debug {
		l.logger.Printf(l.prefix+": "+f, args...)
	}
}

func (l *Logger) Errorf(f string, args ...interface{}) error {
	return fmt.Errorf(l.prefix+": "+f, args...)
}
