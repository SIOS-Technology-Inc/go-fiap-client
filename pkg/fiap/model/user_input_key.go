package model

import "time"

type UserInputKey struct {
	ID              string
	Eq              time.Time
	Neq             time.Time
	Lt              time.Time
	Gt              time.Time
	Lteq            time.Time
	Gteq            time.Time
	MinMaxIndicator string
}
