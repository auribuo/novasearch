package types

import (
	"encoding/json"
	"github.com/auribuo/novasearch/types/coordinates"
	"math"
	"time"
)

const (
	MagnitudeUgc2 = 17.0
	DistanceUgc2  = 268.0
	BaseTimeUgc2  = 536
)

const readoutTime = 3

type Galaxy struct {
	Morphology            string                            `json:"morphology"`
	Id                    int                               `json:"id"`
	PreferredName         string                            `json:"preferredName"`
	Magnitude             float64                           `json:"magnitude"`
	EquatorialCoordinates coordinates.EquatorialCoordinates `json:"equatorialCoordinates"`
	AzimuthalCoordinates  coordinates.AzimuthalCoordinates  `json:"azimuthalCoordinates"`
	SemiMajorAxis         float64                           `json:"semiMajorAxis"`
	SemiMinorAxis         float64                           `json:"semiMinorAxis"`
	Redshift              float64                           `json:"redshift"`
	visited               int
	timestamp             time.Time
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) Mark(t time.Time) {
	g.timestamp = t
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) At() time.Time {
	return g.timestamp
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) String() string {
	jsonString, _ := json.Marshal(g)
	return string(jsonString)
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.Extend())
}

//goland:noinspection GoMixedReceiverTypes
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

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) Visit() {
	g.visited = 1
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) Reset() {
	g.visited = 0
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) Distance() float64 {
	return g.Redshift * 299792.458 / 70
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) Quality() float64 {
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

	if (g.AzimuthalCoordinates != coordinates.AzimuthalCoordinates{}) {
		height := g.AzimuthalCoordinates.Elevation
		heightWeight := math.Pow((height-30)/60, 1/3)
		quality *= heightWeight
	}

	return quality * float64(g.visited)
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) calculateExposure(referenceTime float64) float64 {
	var distance float64
	if g.Distance() > 0 {
		distance = g.Distance()
	} else {
		distance = 100
	}
	return referenceTime * math.Pow(distance/DistanceUgc2, 2)
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) slewFunction(distance float64) float64 {
	return 1/2.0*distance + 6
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) DistanceTo(other Ratable) float64 {
	return g.AzimuthalCoordinates.DistanceTo(other.Position())
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) WaitTime(distance float64) int {
	slewTime := g.slewFunction(distance)
	balanceTime := slewTime / 1.5
	exposureTime := g.calculateExposure(BaseTimeUgc2)
	exposureCheck := exposureTime / 2

	return int(math.Ceil(slewTime + balanceTime + exposureTime + readoutTime + exposureCheck + readoutTime))
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) Position() coordinates.AzimuthalCoordinates {
	return g.AzimuthalCoordinates
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) Exposure() float64 {
	return g.calculateExposure(BaseTimeUgc2)
}

//goland:noinspection GoMixedReceiverTypes
func (g Galaxy) Extend() ExtendedGalaxy {
	return ExtendedGalaxy{
		Galaxy:   g,
		Distance: g.Distance(),
	}
}
