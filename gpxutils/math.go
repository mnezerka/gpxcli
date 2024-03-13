package gpxutils

import (
	"math"
)

func RoundDecimals(n float64, decimals int) float64 {
	k := math.Pow(10, float64(decimals))
	return math.Round(n*k) / k
}
