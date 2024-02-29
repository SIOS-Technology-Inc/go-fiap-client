package fiap

import (
	"time"
)

func TimeTostring(t time.Time) string {
	return t.Format(time.RFC3339)
}

func CheckAndConvertTime(in time.Time) string {
	if in.IsZero() {
		return ""
	} else {
		return TimeTostring(in)
	}
}