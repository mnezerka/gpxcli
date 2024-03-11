package gpxutils

import "github.com/tkrajina/gpxgo/gpx"

func SliceInsertAtStart(s []*gpx.GPXPoint, p *gpx.GPXPoint) []*gpx.GPXPoint {
	s = append(s, nil) // Step 1 - allocate space for one more element
	copy(s[1:], s[0:]) // Step 2 - compy existing content, shifted by one
	s[0] = p           // Step 3 - insert new element
	return s
}
