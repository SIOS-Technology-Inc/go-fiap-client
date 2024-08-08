package tools

import (
	"time"
)

// この関数の説明を追加する
func TimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
