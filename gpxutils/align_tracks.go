/*
 * blog: https://steemit.com/programming/@bitcalm/how-to-compare-gps-tracks
 * impl: https://github.com/jonblack/cmpgpx/tree/master/examples
 */
package gpxutils

import (
	"math"

	"github.com/apex/log"
	"github.com/tkrajina/gpxgo/gpx"
)

func similarity(p1, p2 gpx.GPXPoint) float64 {

	d := HarversineDistance(p1, p2)

	return -d
}

// Needleman-Wunsch algorithm adapted for gps tracks.
func AlignTracks(track1, track2 *gpx.GPX, gap_penalty float64) ([]*gpx.GPXPoint, []*gpx.GPXPoint, error) {

	l := log.WithFields(log.Fields{
		"comp": "gpxutils/AlignTracks",
	})

	l.Info("Aligning tracks")

	a1 := []*gpx.GPXPoint{}
	a2 := []*gpx.GPXPoint{}

	points1, err := GpxFileToPoints(track1)
	if err != nil {
		return a1, a2, err
	}

	points2, err := GpxFileToPoints(track2)
	if err != nil {
		return a1, a2, err
	}

	l.Infof("  track1 points: %d", len(points1))
	l.Infof("  track2 points: %d", len(points2))

	// construct f-matrix and fill it with zeros
	f := make([][]float64, len(points1))
	for i := 0; i < len(points1); i++ {
		f[i] = make([]float64, len(points2))
		for j := 0; j < len(points2); j++ {
			f[i][j] = 0
		}
	}

	for i := 0; i < len(points1); i++ {
		f[i][0] = gap_penalty * float64(i)
	}

	for j := 0; j < len(points2); j++ {
		f[0][j] = gap_penalty * float64(j)
	}

	for i := 1; i < len(points1); i++ {
		t1 := points1[i]
		for j := 1; j < len(points2); j++ {
			t2 := points2[j]
			match := f[i-1][j-1] + similarity(t1, t2)
			delete := f[i-1][j] + gap_penalty
			insert := f[i][j-1] + gap_penalty
			f[i][j] = math.Max(match, math.Max(delete, insert))
		}
	}

	// backtrack to create alignment
	i := len(points1) - 1
	j := len(points2) - 1

	for i > 0 || j > 0 {
		if i > 0 && j > 0 && f[i][j] == f[i-1][j-1]+similarity(points1[i], points2[j]) {
			a1 = sliceInsertAtStart(a1, &points1[i])
			a2 = sliceInsertAtStart(a2, &points2[j])
			i -= 1
			j -= 1
		} else if i > 0 && f[i][j] == f[i-1][j]+gap_penalty {
			a1 = sliceInsertAtStart(a1, &points1[i])
			a2 = sliceInsertAtStart(a2, nil)
			i -= 1
		} else if j > 0 && f[i][j] == f[i][j-1]+gap_penalty {
			a1 = sliceInsertAtStart(a1, nil)
			a2 = sliceInsertAtStart(a2, &points2[j])
			j -= 1
		}
	}

	l.Info("After alignment")
	l.Infof("  track1 points: %d", len(a1))
	l.Infof("  track2 points: %d", len(a2))

	return a1, a2, nil
}

func sliceInsertAtStart(s []*gpx.GPXPoint, p *gpx.GPXPoint) []*gpx.GPXPoint {
	s = append([]*gpx.GPXPoint{p}, s...)
	return s
}
