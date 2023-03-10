package coordinates

type EquatorialCoordinates struct {
	RightAscension float64 `json:"rightAscension"`
	Declination    float64 `json:"declination"`
}

func RightAscensionToDegrees(hours float64, minutes float64, seconds float64) float64 {
	return hours + minutes/60 + seconds/3600
}

func DeclinationToDegrees(degrees float64, minutes float64, seconds float64) float64 {
	return (degrees*3600 + minutes*60 + seconds) * 360 / 86400
}
