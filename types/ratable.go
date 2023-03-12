package types

type Ratable interface {
	Visitable
	Timestamped
	DistanceTo(other Ratable) float64
	Quality() float64
	Position() AzimuthalCoordinates
	WaitTime(distance float64) int
	Exposure() float64
}
