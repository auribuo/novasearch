package types

import (
	"github.com/auribuo/novasearch/types/coordinates"
)

type Ratable interface {
	Visitable
	Timestamped
	DistanceTo(other Ratable) float64
	Quality() float64
	Position() coordinates.AzimuthalCoordinates
	WaitTime(distance float64) int
	Exposure() float64
}
