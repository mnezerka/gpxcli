package gpxutils_test

import (
	"mnezerka/gpxcli/gpxutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkrajina/gpxgo/gpx"
)

func TestHaversineDistance(t *testing.T) {

	var tests = []struct {
		p1        gpx.GPXPoint
		p2        gpx.GPXPoint
		outMeters float64
	}{
		{
			gpx.GPXPoint{Point: gpx.Point{Latitude: 22.55, Longitude: 43.12}},  // Rio de Janeiro, Brazil
			gpx.GPXPoint{Point: gpx.Point{Latitude: 13.45, Longitude: 100.28}}, // Bangkok, Thailand
			6101393.716,
		},
		{
			gpx.GPXPoint{Point: gpx.Point{Latitude: 20.10, Longitude: 57.30}}, // Port Louis, Mauritius
			gpx.GPXPoint{Point: gpx.Point{Latitude: 0.57, Longitude: 100.21}}, // Padang, Indonesia
			5151308.531,
		},
		{
			gpx.GPXPoint{Point: gpx.Point{Latitude: 45.04, Longitude: 7.42}},  // Turin, Italy
			gpx.GPXPoint{Point: gpx.Point{Latitude: 3.09, Longitude: 101.42}}, // Kuala Lumpur, Malaysia
			10089438.164,
		},
	}

	for _, input := range tests {
		assert.Equal(t, input.outMeters, gpxutils.RoundDecimals(gpxutils.HarversineDistance(input.p1, input.p2), 3))
	}
}

func TestBearing(t *testing.T) {
	// single point
	p1 := gpx.GPXPoint{Point: gpx.Point{Latitude: 10.0, Longitude: 10.0}}
	assert.Equal(t, 0.0, gpxutils.RoundDecimals(gpxutils.Bearing(p1, p1), 1))

	// poinst above each other -> 0
	p2 := gpx.GPXPoint{Point: gpx.Point{Latitude: 20.0, Longitude: 10.0}}
	assert.Equal(t, 0.0, gpxutils.RoundDecimals(gpxutils.Bearing(p1, p2), 1))

	// points side by side equally high
	p2 = gpx.GPXPoint{Point: gpx.Point{Latitude: 10.0, Longitude: 20.0}}
	assert.Equal(t, 89.1, gpxutils.RoundDecimals(gpxutils.Bearing(p1, p2), 1))

	// points in 45 grade direction
	p2 = gpx.GPXPoint{Point: gpx.Point{Latitude: 20.0, Longitude: 20.0}}
	assert.Equal(t, 42.8, gpxutils.RoundDecimals(gpxutils.Bearing(p1, p2), 1))
}
