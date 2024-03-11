package gpxutils_test

import (
	"math"
	"mnezerka/gpxcli/gpxutils"
	"testing"

	"github.com/tkrajina/gpxgo/gpx"
)

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

func TestHaversineDistance(t *testing.T) {
	for _, input := range tests {
		meters := math.Round(gpxutils.HarversineDistance(input.p1, input.p2)*1000) / 1000

		if input.outMeters != meters {
			t.Errorf("fail: want %f got %f", input.outMeters, meters)
		}
	}
}
