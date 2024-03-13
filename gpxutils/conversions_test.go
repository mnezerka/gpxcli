package gpxutils_test

import (
	"math"
	"mnezerka/gpxcli/gpxutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRadiansToDegrees(t *testing.T) {
	// pi (3.14..) radians is equal to 180 degrees, which means that 1 radian is equal to 57.2957795 degrees.
	assert.Equal(t, 57.296, gpxutils.RoundDecimals(gpxutils.RadiansToDegrees(1), 3))
	assert.Equal(t, 180.0, gpxutils.RoundDecimals(gpxutils.RadiansToDegrees(math.Pi), 3))
}

func TestDegreesToRadians(t *testing.T) {
	assert.Equal(t, 3.142, gpxutils.RoundDecimals(gpxutils.DegreesToRadians(180), 3))

}
