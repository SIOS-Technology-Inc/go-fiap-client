package tools

import (
	"log"
)

var DEBUG = true

func DebugLogPrintln(a ...interface{}) {
	if DEBUG {
		log.Println(a...)
	}
}

func DebugLogPrintf(format string, a ...interface{}) {
	if DEBUG {
		log.Printf(format, a...)
	}
}
