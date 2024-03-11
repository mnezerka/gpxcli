package gpxutils

import (
	"math"

	"github.com/tkrajina/gpxgo/gpx"
)

// degreesToRadians converts from degrees to radians.
func DegreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// Join all the points from all tracks and segments for the track into a single list
// func GpxToPoints(gpx *gpx.GPX) (error, []gpx.GPXPoint) {
func GpxFileToPoints(gpxFile *gpx.GPX) ([]gpx.GPXPoint, error) {

	var result []gpx.GPXPoint

	for _, track := range gpxFile.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				result = append(result, point)
			}
		}
	}
	return result, nil
}
