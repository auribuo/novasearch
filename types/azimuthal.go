package types

import (
	"math"

	"github.com/auribuo/novasearch/util"
)

type AzimuthalCoordinates struct {
	Azimuth   float64 `json:"azimuth"`
	Elevation float64 `json:"elevation"`
}

func (coordinate AzimuthalCoordinates) DistanceTo(other AzimuthalCoordinates) float64 {
	deltaAzimuth, deltaElevation := coordinate.Azimuth-other.Azimuth, coordinate.Elevation-other.Elevation
	distanceRad := math.Acos(math.Cos(deltaAzimuth) * math.Cos(deltaElevation))
	return util.ToDegrees(distanceRad)
}
