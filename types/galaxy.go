package types

import (
	"encoding/json"
	"math"
	"time"

	"github.com/auribuo/novasearch/sql"
)

const (
	MagnitudeUgc2 = 17.0
	DistanceUgc2  = 268.0
	BaseTimeUgc2  = 536
)

const readoutTime = 3

type extendedGalaxy Galaxy

type Galaxy struct {
	Morphology            string                `json:"morphology"`
	Id                    int                   `json:"id"`
	PreferredName         string                `json:"preferredName"`
	Magnitude             float64               `json:"magnitude"`
	EquatorialCoordinates EquatorialCoordinates `json:"equatorialCoordinates"`
	AzimuthalCoordinates  AzimuthalCoordinates  `json:"azimuthalCoordinates"`
	SemiMajorAxis         float64               `json:"semiMajorAxis"`
	SemiMinorAxis         float64               `json:"semiMinorAxis"`
	Redshift              float64               `json:"redshift"`
	visited               bool
	timestamp             time.Time
}

func (g *Galaxy) Mark(t time.Time) {
	g.timestamp = t
}

func (g *Galaxy) At() time.Time {
	return g.timestamp
}

func (g *Galaxy) String() string {
	jsonString, _ := json.Marshal(g)
	return string(jsonString)
}

func (g *Galaxy) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		*extendedGalaxy
		Distance float64 `json:"distance"`
		Quality  float64 `json:"quality"`
		Visited  bool    `json:"visited"`
	}{
		extendedGalaxy: (*extendedGalaxy)(g),
		Distance:       g.Distance(),
		Quality:        g.Quality(),
		Visited:        g.visited,
	})
}

func (g *Galaxy) UnmarshalJSON(data []byte) error {
	type Alias Galaxy
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(g),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

func (g *Galaxy) Visit() {
	g.visited = true
}

func (g *Galaxy) Reset() {
	g.visited = false
}

func (g *Galaxy) Distance() float64 {
	return g.Redshift * 299792.458 / 70
}

func (g *Galaxy) Quality() float64 {
	if g.visited {
		return 0
	}

	var typeWeight float64

	switch g.Morphology {
	case "E":
		typeWeight = 0.09984127
	case "S0":
		typeWeight = 0.08793651
	case "Sa":
		typeWeight = 0.13253968
	case "Sb":
		typeWeight = 0.14587302
	case "Sbc":
		typeWeight = 0.16936508
	case "Sc":
		typeWeight = 0.21619048
	case "Scd":
		typeWeight = 0.19746032
	default:
		typeWeight = 0.149887
	}

	quality := typeWeight * math.Pow(10, (g.Magnitude-MagnitudeUgc2)/-2.5)

	if (g.AzimuthalCoordinates != AzimuthalCoordinates{}) {
		height := g.AzimuthalCoordinates.Elevation
		heightWeight := math.Pow((height-30)/60, 1/3.0)
		quality *= heightWeight
	}

	if g.visited {
		return 0
	}

	return quality
}

func (g *Galaxy) calculateExposure(referenceTime float64) float64 {
	var distance float64
	if g.Distance() > 0 {
		distance = g.Distance()
	} else {
		distance = 100
	}
	return referenceTime * math.Pow(distance/DistanceUgc2, 2)
}

func (g *Galaxy) slewFunction(distance float64) float64 {
	return 1/2.0*distance + 6
}

func (g *Galaxy) DistanceTo(other Ratable) float64 {
	return g.AzimuthalCoordinates.DistanceTo(other.Position())
}

func (g *Galaxy) WaitTime(distance float64) int {
	slewTime := g.slewFunction(distance)
	balanceTime := slewTime / 1.5
	exposureTime := g.calculateExposure(BaseTimeUgc2)
	exposureCheck := exposureTime / 2

	return int(math.Ceil(slewTime + balanceTime + exposureTime + readoutTime + exposureCheck + readoutTime))
}

func (g *Galaxy) Position() AzimuthalCoordinates {
	return g.AzimuthalCoordinates
}

func (g *Galaxy) Exposure() float64 {
	return g.calculateExposure(BaseTimeUgc2)
}

func CalculateViewports(source []Galaxy, fov Fov, location Location, dateTime time.Time, gridApprox float64) []Viewport {
	var viewports = make(map[Tuple[float64, float64]]Viewport)
	xStep := fov.Width * gridApprox
	yStep := fov.Height * gridApprox

	for _, galaxy := range source {
		xApprox := nearestDegree(galaxy.EquatorialCoordinates.RightAscension, xStep)
		yApprox := nearestDegree(galaxy.EquatorialCoordinates.Declination, yStep)

		xyTuple := NewTuple(xApprox, yApprox)

		if viewport, ok := viewports[xyTuple]; ok {
			viewport.Galaxies = append(viewport.Galaxies, galaxy)
			viewports[xyTuple] = viewport
		} else {
			viewport = Viewport{
				EquatorialPosition: EquatorialCoordinates{
					RightAscension: xApprox,
					Declination:    yApprox,
				},
				Galaxies: []Galaxy{galaxy},
			}

			if location != (Location{}) && dateTime != (time.Time{}) {
				topLeft := EquatorialCoordinates{
					RightAscension: viewport.EquatorialPosition.RightAscension - fov.Width,
					Declination:    viewport.EquatorialPosition.Declination + fov.Height,
				}
				topRight := EquatorialCoordinates{
					RightAscension: viewport.EquatorialPosition.RightAscension + fov.Width,
					Declination:    viewport.EquatorialPosition.Declination + fov.Height,
				}
				bottomLeft := EquatorialCoordinates{
					RightAscension: viewport.EquatorialPosition.RightAscension - fov.Width,
					Declination:    viewport.EquatorialPosition.Declination - fov.Height,
				}
				bottomRight := EquatorialCoordinates{
					RightAscension: viewport.EquatorialPosition.RightAscension + fov.Width,
					Declination:    viewport.EquatorialPosition.Declination - fov.Height,
				}
				viewport.ViewportEdges = ViewportEdges{
					TopLeft:     topLeft.ToAzimuthalCoordinates(dateTime, location),
					TopRight:    topRight.ToAzimuthalCoordinates(dateTime, location),
					BottomLeft:  bottomLeft.ToAzimuthalCoordinates(dateTime, location),
					BottomRight: bottomRight.ToAzimuthalCoordinates(dateTime, location),
				}
				viewport.AzimuthalPosition = viewport.EquatorialPosition.ToAzimuthalCoordinates(dateTime, location)
			}

			viewports[xyTuple] = viewport
		}
	}
	return sql.Values(viewports)
}

func nearestDegree(degree float64, step float64) float64 {
	var nearestDeg float64

	n := degree / step

	if math.Abs(degree-n*step) < math.Abs(degree-(n+1)*step) {
		nearestDeg = n * step
	} else {
		nearestDeg = (n + 1) * step
	}

	return nearestDeg
}
