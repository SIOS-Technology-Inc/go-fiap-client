package tools

import (
	"log"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = 0
	LogLevelError  LogLevel = 3
)

func SetLogLevel(level string) {
	switch level {
	case "debug":
		LOGLEVEL = LogLevelDebug
	case "error":
		LOGLEVEL = LogLevelError
	default:
		log.Printf("Invalid log level: %s. Set string 'debug' or 'error'. Setting to default 'error'.", level)
		LOGLEVEL = LogLevelError
	}
}

var LOGLEVEL LogLevel = LogLevelError


func DebugLogPrintf(format string, a ...interface{}) {
	logPrintf(LogLevelDebug, format, a...)
}

func ErrorLogPrintf(format string, a ...interface{}) {
	logPrintf(LogLevelError, format, a...)
}	


func logPrintf(l LogLevel, format string, a ...interface{}) {
	if l >= LOGLEVEL {
		log.Printf(logLevelToString(l) + ": " + format, a...)
	}
}

func logLevelToString(l LogLevel) string {
	if l == LogLevelDebug {
		return "Debug"
	} else if l == LogLevelError {
		return "Error"
	} else {
		return "Unknown"
	}
}