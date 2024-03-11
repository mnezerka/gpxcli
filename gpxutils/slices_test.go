package gpxutils_test

import (
	"mnezerka/gpxcli/gpxutils"
	"reflect"
	"testing"

	"github.com/tkrajina/gpxgo/gpx"
)

func TestSliceInsertAtStart(t *testing.T) {

	s := []*gpx.GPXPoint{}

	s = gpxutils.SliceInsertAtStart(s, nil)

	expected := []*gpx.GPXPoint{nil}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("fail: want %v got %v", expected, s)
	}

	s = gpxutils.SliceInsertAtStart(s, nil)

	expected = []*gpx.GPXPoint{nil, nil}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("fail: want %v got %v", expected, s)
	}

	point := gpx.GPXPoint{Point: gpx.Point{Latitude: 23, Longitude: 58}}
	point2 := gpx.GPXPoint{Point: gpx.Point{Latitude: 11, Longitude: 22}}
	s = gpxutils.SliceInsertAtStart(s, &point)
	s = gpxutils.SliceInsertAtStart(s, &point2)

	expected = []*gpx.GPXPoint{&point2, &point, nil, nil}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("fail: want %v got %v", expected, s)
	}
}
