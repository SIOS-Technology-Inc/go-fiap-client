package tools

import (
	"log"
)

type logLevel int

const (
	LogLevelDebug logLevel = 0
	LogLevelError  logLevel = 3
)

var setLogLevel logLevel = LogLevelError

func SetLogLevel(l logLevel) {
	setLogLevel = l
}

func LogPrintf(printLogLevel logLevel, format string, a ...interface{}) {
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