package tools

import (
	"time"
)

func TimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	} else {
		return t.Format(time.RFC3339)
	}
}