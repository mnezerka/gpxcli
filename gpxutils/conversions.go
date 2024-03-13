package gpxutils

import (
	"fmt"
	"math"

	"github.com/apex/log"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/tkrajina/gpxgo/gpx"
)

// degreesToRadians converts from degrees to radians.
func DegreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// converts an angle from radians to degrees.
func RadiansToDegrees(r float64) float64 {
	return r * 180 / math.Pi
}

// Join all the points from all tracks and segments for the track into a single list
// func GpxToPoints(gpx *gpx.GPX) (error, []gpx.GPXPoint) {
func GpxFileToPoints(gpxFile *gpx.GPX) ([]gpx.GPXPoint, error) {

	var result []gpx.GPXPoint

	for _, track := range gpxFile.Tracks {
		for _, segment := range track.Segments {
			result = append(result, segment.Points...)
		}
	}
	return result, nil
}

func RenderTracksDataToGeoJson(tracks []RenderTrackData, points bool) ([]byte, error) {

	l := log.WithFields(log.Fields{
		"comp": "gpxutils/RenderTracksDataToGeoJson",
	})

	l.Infof("Conversion []RenderTrackData -> GeoJson, tracks: %d, draw-points: %v", len(tracks), points)

	fc := geojson.NewFeatureCollection()

	for i := 0; i < len(tracks); i++ {

		track := tracks[i]
		l.Infof("  track: %d, track-points: %d", i, len(track.Points))

		line := orb.LineString{}
		for j := 0; j < len(track.Points); j++ {
			line = append(line, orb.Point{track.Points[j].Longitude, track.Points[j].Latitude})
		}

		feature := geojson.NewFeature(line)
		feature.Properties["color"] = fmt.Sprintf("#%02x%02x%02x", track.Color.R, track.Color.G, track.Color.B)

		fc.Append(feature)

		// if points are enabled, add
		if points {
			for j := 0; j < len(track.Points); j++ {
				pointFeature := geojson.NewFeature(orb.Point{track.Points[j].Longitude, track.Points[j].Latitude})
				pointFeature.Properties["color"] = fmt.Sprintf("#%02x%02x%02x", track.Color.R, track.Color.G, track.Color.B)
				fc.Append(pointFeature)
			}
		}
	}

	rawJson, err := fc.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}

	return rawJson, nil
}

func GpxPointsToGeoJson(points []gpx.GPXPoint) ([]byte, error) {

	fc := geojson.NewFeatureCollection()
	for i := 0; i < len(points); i++ {
		fc.Append(geojson.NewFeature(orb.Point{points[i].Point.Latitude, points[i].Point.Longitude}))
	}
	rawJson, err := fc.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}

	return rawJson, nil
}

func Reverse(points []gpx.GPXPoint) ([]gpx.GPXPoint, error) {
	for i, j := 0, len(points)-1; i < j; i, j = i+1, j-1 {
		points[i], points[j] = points[j], points[i]
	}

	return points, nil
}
