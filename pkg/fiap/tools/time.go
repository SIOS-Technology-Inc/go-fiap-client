package tools

import (
	"time"
)

/*
TimeToString returns the string representation(RFC3339) of the given time.Time.

TimeToStringは指定されたtime.Timeの文字列表現(RFC3339)を返します。
*/
func TimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
