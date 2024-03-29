package model

import "time"

type ProcessedPoint struct {
	Times  []time.Time `json:"time"`
	Values []string    `json:"value"`
}

type ProcessedPointSet struct {
	PointSetID []string `json:"point_set_id"`
	PointID    []string `json:"point_ids"`
}