package model

import "time"

type ProcessedValue struct {
	Time  time.Time `json:"time"`
	Value string    `json:"value"`
}