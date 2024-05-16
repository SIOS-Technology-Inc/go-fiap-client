package tools

import (
	"log"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = 0
	LogLevelError  LogLevel = 3
)

var setLogLevel LogLevel = LogLevelError

func SetLogLevel(l LogLevel) {
	setLogLevel = l
}

func LogPrintf(printLogLevel LogLevel, format string, a ...interface{}) {
	if printLogLevel >= setLogLevel {
		var levelStr string
		if printLogLevel == LogLevelDebug {
			levelStr = "Debug"
		} else if printLogLevel == LogLevelError {
			levelStr = "Error"
		}
		log.Printf(levelStr + ": " + format, a...)
	}
}