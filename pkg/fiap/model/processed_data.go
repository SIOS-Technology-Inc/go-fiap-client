package model

import "time"

type ProcessedValue struct {
	Time  time.Time `json:"time"`
	Value string    `json:"value"`
}

type ProcessedPoint struct {
	Values []ProcessedValue `json:"values"`
}

type ProcessedPointSet struct {
	PointSetID []string `json:"point_set_id"`
	PointID    []string `json:"point_ids"`
}