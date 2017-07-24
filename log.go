package summer

import (
	"log"
	"os"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var Logger = NewSimpleLogger(InfoLevel)
var logger = Logger.Module("summer")

type SimpleLogger struct {
	level LogLevel
}

func NewSimpleLogger(logLevel LogLevel) *SimpleLogger {
	return &SimpleLogger{logLevel}
}
func (logger *SimpleLogger) Module(module string) *SimpleLog {
	return &SimpleLog{SimpleLogger: logger, log: log.New(os.Stderr, "["+module+"]", log.LstdFlags)}
}

type SimpleLog struct {
	*SimpleLogger
	log *log.Logger
}

func NewSimpleLog(module string, logLevel LogLevel) *SimpleLog {
	return NewSimpleLogger(logLevel).Module(module)
}

func SetLogLevel(logLevel LogLevel) {
	Logger.level = logLevel
}
func (log *SimpleLog) Debug(args ...interface{}) {
	if DebugLevel < log.level {
		return
	}
	log.log.Println(args...)
}
func (log *SimpleLog) Error(args ...interface{}) {
	if ErrorLevel < log.level {
		return
	}
	log.log.Println(args...)
}
func (log *SimpleLog) Println(args ...interface{}) {
	log.log.Println(args...)
}
func (log *SimpleLog) Warn(args ...interface{}) {
	if WarnLevel < log.level {
		return
	}
	log.log.Println(args...)
}

func (log *SimpleLog) Panic(args ...interface{}) {
	if PanicLevel < log.level {
		return
	}
	log.log.Panicln(args...)
}

func (log *SimpleLog) Fatal(args ...interface{}) {
	if FatalLevel < log.level {
		return
	}
	log.log.Fatalln(args...)
}

func (log *SimpleLog) Info(args ...interface{}) {
	if InfoLevel < log.level {
		return
	}
	log.log.Println(args...)
}
