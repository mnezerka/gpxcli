package gpxutils

import (
	"image/color"

	sm "github.com/flopp/go-staticmaps"

	"github.com/apex/log"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/tkrajina/gpxgo/gpx"
)

var ColorPalette = [...]color.Color{
	color.RGBA{0xff, 0, 0, 0xff},
	color.RGBA{0, 0xff, 0, 0xff},
	color.RGBA{0, 0, 0xff, 0xff},
}

type RenderTrackData struct {
	Points []gpx.GPXPoint
	Color  color.Color
}

func GetColorForIndex(i int) color.Color {
	i = i % len(ColorPalette)
	return ColorPalette[i]
}

func RenderTracks(tracks []RenderTrackData) error {

	l := log.WithFields(log.Fields{
		"comp": "gpxutils/render",
	})

	ctx := sm.NewContext()
	ctx.SetSize(800, 600)

	for i := 0; i < len(tracks); i++ {
		track := tracks[i]
		for j := 0; j < len(track.Points); j++ {
			point := track.Points[j]

			l.Debugf("rendering %v", point)

			ctx.AddObject(
				sm.NewMarker(
					s2.LatLngFromDegrees(point.Latitude, point.Longitude),
					track.Color,
					16.0,
				),
			)
		}
	}

	img, err := ctx.Render()
	if err != nil {
		return err
	}

	if err := gg.SavePNG("my-map.png", img); err != nil {
		return err
	}

	return nil
}
