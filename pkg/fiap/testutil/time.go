package testutil

import "time"

// TimeToTimep returns a pointer to the given time.Time.
// TimeToTimepは指定されたtime.Timeのポインタを返します。
func TimeToTimep(t time.Time) *time.Time {
	return &t
}