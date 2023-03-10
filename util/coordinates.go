package util

import "math"

func ToDegrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

func Delta(x float64, y float64) float64 {
	return math.Abs(x - y)
}
