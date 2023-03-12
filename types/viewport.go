package types

import (
	"math"
	"time"

	"github.com/auribuo/novasearch/util"
)

type ViewportEdges struct {
	TopLeft     AzimuthalCoordinates `json:"topLeft"`
	TopRight    AzimuthalCoordinates `json:"topRight"`
	BottomLeft  AzimuthalCoordinates `json:"bottomLeft"`
	BottomRight AzimuthalCoordinates `json:"bottomRight"`
}

type Viewport struct {
	EquatorialPosition EquatorialCoordinates `json:"equatorialPosition"`
	AzimuthalPosition  AzimuthalCoordinates  `json:"azimuthalPosition"`
	Galaxies           []Galaxy              `json:"galaxies"`
	ViewportEdges      ViewportEdges         `json:"viewportEdges"`
}

func (v Viewport) At() time.Time {
	if (v.Galaxies == nil) || (len(v.Galaxies) == 0) {
		return time.Time{}
	}
	galaxy := v.Galaxies[0]
	return galaxy.At()
}

func (v Viewport) DistanceTo(other Ratable) float64 {
	center := AzimuthalCoordinates{
		Elevation: v.ViewportEdges.BottomLeft.Elevation + util.Delta(v.ViewportEdges.TopLeft.Elevation, v.ViewportEdges.BottomLeft.Elevation)/2,
		Azimuth:   v.ViewportEdges.BottomLeft.Azimuth + util.Delta(v.ViewportEdges.BottomRight.Azimuth, v.ViewportEdges.BottomLeft.Azimuth)/2,
	}
	deltaElevation := util.Delta(center.Elevation, other.Position().Elevation)
	deltaAzimuth := util.Delta(center.Azimuth, other.Position().Azimuth)

	return math.Acos(math.Cos(deltaElevation) * math.Cos(deltaAzimuth))
}

func (v Viewport) Quality() float64 {
	sum := 0.0
	for _, galaxy := range v.Galaxies {
		sum += galaxy.Quality()
	}
	return sum
}

func (v Viewport) Position() AzimuthalCoordinates {
	return v.AzimuthalPosition
}

func (v Viewport) WaitTime(distance float64) int {
	max := math.MinInt
	for _, galaxy := range v.Galaxies {
		if galaxy.WaitTime(distance) > max {
			max = galaxy.WaitTime(distance)
		}
	}
	return max
}

func (v Viewport) Exposure() float64 {
	max := math.SmallestNonzeroFloat64
	for _, galaxy := range v.Galaxies {
		if galaxy.Exposure() > max {
			max = galaxy.Exposure()
		}
	}
	return max
}

func (v Viewport) Visit() {
	for _, galaxy := range v.Galaxies {
		galaxy.Visit()
	}
}

func (v Viewport) Reset() {
	for _, galaxy := range v.Galaxies {
		galaxy.Reset()
	}
}

func (v Viewport) Mark(t time.Time) {
	for _, galaxy := range v.Galaxies {
		galaxy.Mark(t)
	}
}
