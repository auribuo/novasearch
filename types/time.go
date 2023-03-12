package types

import "time"

type TimeStat struct {
	TotalTime time.Duration `json:"total_time"`
	UsedTime  time.Duration `json:"used_time"`
	Diff      time.Duration `json:"diff_time"`
}
