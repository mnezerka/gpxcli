package gpxutils

import (
	"math"

	"github.com/tkrajina/gpxgo/gpx"
)

const (
	earthRaidusMeters = 6378160 // radius of the earth in meters.
)

// Distance calculates the shortest path between two coordinates on the surface
// of the Earth in meters
func HarversineDistance(p, q gpx.GPXPoint) float64 {
	lat1 := DegreesToRadians(p.Latitude)
	lon1 := DegreesToRadians(p.Longitude)
	lat2 := DegreesToRadians(q.Latitude)
	lon2 := DegreesToRadians(q.Longitude)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * earthRaidusMeters
}
