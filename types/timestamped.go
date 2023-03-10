package types

import "time"

type Timestamped interface {
	Mark(t time.Time)
	At() time.Time
}
