package gpxutils

import (
	"math"

	"github.com/apex/log"
	"github.com/tkrajina/gpxgo/gpx"
)

const (
	// radius of the earth in meters.
	earthRaidusMeters = 6378160

	// one degree in meters:
	oneDegree = (2 * math.Pi * earthRaidusMeters) / 360 // 111.319 km
)

// Distance calculates the shortest path between two coordinates on the surface
// of the Earth in meters
// available also here: https://github.com/tkrajina/gpxgo/blob/master/gpx/geo.go#L55
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

/*
 * Calculates the initial bearing between point1 and point2 relative to north
 * (zero degrees).
 */
func Bearing(point1, point2 gpx.GPXPoint) float64 {

	lat1r := DegreesToRadians(point1.Latitude)
	lat2r := DegreesToRadians(point2.Latitude)
	dlon := DegreesToRadians(point2.Longitude - point1.Longitude)

	y := math.Sin(dlon) * math.Cos(lat2r)
	x := math.Cos(lat1r)*math.Sin(lat2r) - math.Sin(lat1r)*math.Cos(lat2r)*math.Cos(dlon)
	return RadiansToDegrees(math.Atan2(y, x))
}

/*
 * Interpolates points so that the distance between each point is equal
 * to `distance` in meters.
 *
 * Only latitude and longitude are interpolated; time and elavation are not
 * interpolated and should not be relied upon.
 *
 * TODO: Interpolate elevation and time.
 */
func InterpolateDistance(points []gpx.GPXPoint, distance float64) ([]gpx.GPXPoint, error) {

	l := log.WithFields(log.Fields{
		"comp": "gpxutils/InterpolateDistance",
	})

	l.Infof("Distributing points evenly every %f meters", distance)
	l.Infof("  points-count-initial: %d", len(points))

	var result []gpx.GPXPoint

	if len(points) == 0 {
		return result, nil
	}

	var d float64 = 0
	var i int = 0
	var p1, p2 gpx.GPXPoint

	for i < len(points) {
		// first point doesn't need any processing
		if i == 0 {
			result = append(result, points[0])
			i += 1
			continue
		}

		if d == 0 {
			p1 = result[len(result)-1]
		} else {
			p1 = points[i-1]
		}

		p2 = points[i]

		d += HarversineDistance(p1, p2)

		if d >= distance {
			bearing := Bearing(p1, p2)
			p2_copy := MoveByAngleAndDistance2(p2, bearing, -(d - distance))
			result = append(result, p2_copy)
			l.Debugf("  adding new point at lat: %f lng: %f", p2_copy.Latitude, p2_copy.Longitude)
			d = 0
		} else {
			i += 1
		}
	}
	result = append(result, points[len(points)-1])

	l.Infof("  points-count-final: %d", len(result))

	return result, nil
}

func MoveByAngleAndDistance(p gpx.GPXPoint, angle, distance float64) gpx.GPXPoint {
	coef := math.Cos(DegreesToRadians(p.Latitude))
	verticalDistanceDiff := math.Sin(DegreesToRadians(90-angle)) / oneDegree
	horizontalDistanceDiff := math.Cos(DegreesToRadians(90-angle)) / oneDegree
	latDiff := distance * verticalDistanceDiff
	lonDiff := distance * horizontalDistanceDiff / coef

	result := gpx.GPXPoint{}
	result.Latitude = p.Latitude + latDiff
	result.Longitude = p.Longitude + lonDiff

	return result
}

// angle measured clockwise from due north

func MoveByAngleAndDistance2(p gpx.GPXPoint, angle, distance float64) gpx.GPXPoint {

	angleRad := DegreesToRadians(angle)

	// dx, dy same units as distance
	dx := distance * math.Sin(angleRad)
	dy := distance * math.Cos(angleRad)

	deltaLongitude := dx / (oneDegree * math.Cos(DegreesToRadians(p.Latitude)))
	deltaLatitude := dy / oneDegree

	result := gpx.GPXPoint{}
	result.Latitude = p.Latitude + deltaLatitude
	result.Longitude = p.Longitude + deltaLongitude

	return result
}
