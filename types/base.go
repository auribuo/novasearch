package types

import "time"

type BaseData struct {
	StartTime        time.Time  `json:"startTime"`
	Location         Location   `json:"location"`
	Fov              Fov        `json:"fov"`
	Hemisphere       Hemisphere `json:"hemisphere"`
	MinHeight        float64    `json:"minHeight"`
	MaxSemiMajorAxis float64    `json:"maxSemiMajorAxis"`
	MaxSemiMinorAxis float64    `json:"maxSemiMinorAxis"`
}
