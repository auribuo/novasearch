package ugc

import "time"

type Cache struct {
	LastUpdated time.Time
	Items       []Response
}
