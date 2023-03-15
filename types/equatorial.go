package types

import (
	"math"
	"time"

	"github.com/starryalley/go-julianday"
)

type EquatorialCoordinates struct {
	RightAscension float64 `json:"rightAscension"`
	Declination    float64 `json:"declination"`
}

func RightAscensionToDegrees(hours float64, minutes float64, seconds float64) float64 {
	return (hours*3600 + minutes*60 + seconds) * 360 / 86400
}

func DeclinationToDegrees(degrees float64, minutes float64, seconds float64) float64 {
	if degrees < 0 {
		return degrees - (minutes/60 + seconds/3600)
	}

	return degrees + minutes/60 + seconds/3600
}

func (c *EquatorialCoordinates) ToAzimuthalCoordinates(dateTime time.Time, location Location) AzimuthalCoordinates {
	dateTime = dateTime.UTC()
	longitude := location.Longitude
	latitude := location.Latitude
	rightAscension := math.Mod(c.RightAscension+360.0, 360.0)
	declination := math.Mod(c.Declination+360.0, 360.0)
	julianDate := julianday.Date(dateTime) + 1
	t := (julianDate - 2451545.0) / 36525.0
	gmst0 := 100.46061837 + 36000.770053608*t + 0.000387933*t*t - (t * t * t / 38710000)
	hour := dateTime.Hour()
	absHours := math.Mod(float64((hour)+24.0), 24.0)
	hourInDayFraction := absHours / 24.0
	timeOfDayDeg := hourInDayFraction * 360.0
	timeOfDayRect := float64(timeOfDayDeg) * 1.00273790935
	gmst := math.Mod(gmst0+timeOfDayRect, 360)
	lmst := math.Mod(gmst+longitude, 360)
	hourAngle := math.Mod(lmst-rightAscension+360, 360)
	heightAngle := math.Asin(math.Sin(degToRad(declination))*math.Sin(degToRad(latitude)) +
		math.Cos(degToRad(latitude))*math.Cos(degToRad(declination))*math.Cos(degToRad(hourAngle)))
	azimuth := math.Atan(math.Sin(degToRad(hourAngle)) / (math.Sin(degToRad(latitude))*math.Cos(degToRad(hourAngle)) -
		math.Cos(degToRad(latitude))*math.Tan(degToRad(declination))))
	azimuthDeg := radToDeg(azimuth)
	azimuthDeg = math.Mod(azimuthDeg+360*5, 360)

	if getQuadrant(azimuthDeg) == getQuadrant(hourAngle) {
		azimuthDeg = math.Mod(azimuthDeg+180, 360)
	}

	return AzimuthalCoordinates{
		Azimuth:   azimuthDeg,
		Elevation: radToDeg(heightAngle),
	}
}

func degToRad(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func radToDeg(radians float64) float64 {
	return radians * 180 / math.Pi
}

func getQuadrant(degrees float64) int {
	angle := math.Mod(degrees, 360)
	if angle >= 0 && angle < 90 {
		return 1
	}
	if angle >= 90 && angle < 180 {
		return 2
	}
	if angle >= 180 && angle < 270 {
		return 3
	}
	if angle >= 270 && angle < 360 {
		return 4
	}
	return 1
}
