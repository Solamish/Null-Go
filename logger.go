package nullgo

import (
	"log"
	"os"
)

const (
	LevelTrace = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

var level = LevelTrace

//setter
func SetLevel(l int) {
	level = l
}

//getter
func Level() int {
	return level
}

var NullLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var NullLogger2 = log.New(os.Stdout, "[BULL] ", log.Ldate|log.Ltime)
var NullLogger3 = log.New(os.Stdout, "[BULL-WebSocket] ", log.Ldate|log.Ltime)

func SetLogger(l *log.Logger) {
	NullLogger = l
}

func Trace(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-trace] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-error] "+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-warn] "+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-info] "+format, v...)
	}
}
func Debug(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-debug] "+format, v...)
	}
}

func printError(format string, v ...interface{})  {
	if level <= LevelTrace {
		NullLogger2.Printf(format, v...)
	}
}

func printDebug(format string, v ...interface{})  {
	if level <= LevelTrace {
		NullLogger2.Printf(format, v...)
	}
}

func wsInfo(format string, v ...interface{})  {
	if level <= LevelTrace {
		NullLogger3.Printf(format, v...)
	}
}