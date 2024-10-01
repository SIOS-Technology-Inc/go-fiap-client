package tools

import (
	"log"
)

type logLevel int

/*
LogLevelDebug is a constant of logLevel, representing debug log level.

LogLevelDebugは、デバッグログレベルを表すlogLevelの定数です。
出力ログレベルの設定や、ログprint時のログレベルの指定に使用します。
*/
const LogLevelDebug logLevel = 0
/*
LogLevelError is a constant of logLevel, representing error log level.	

LogLevelErrorは、エラーログレベルを表すlogLevelの定数です。
出力ログレベルの設定や、ログprint時のログレベルの指定に使用します。
*/
const LogLevelError  logLevel = 3

var setLogLevel logLevel = LogLevelError

/*
SetLogLevel sets the output log level.

SetLogLevelは出力ログのレベルを設定します。
*/
func SetLogLevel(l logLevel) {
	setLogLevel = l
}

/*
LogPrintf prints logs with the specified log level.

LogPrintfはログレベルを指定してログのprintを行います。
*/
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