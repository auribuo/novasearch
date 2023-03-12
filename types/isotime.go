package types

import "time"

type IsoTime string

func (t IsoTime) Convert() time.Time {
	result, _ := time.Parse(time.RFC3339, string(t))
	return result
}
