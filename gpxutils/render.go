package gpxutils

import (
	"image/color"
	"os"
	"strings"
	"text/template"

	sm "github.com/flopp/go-staticmaps"

	"github.com/apex/log"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/tkrajina/gpxgo/gpx"
)

var ColorPalette = [...]color.RGBA{
	{0xaf, 0, 0, 0xff},
	{0, 0xaf, 0, 0xff},
	{0, 0, 0xaf, 0xff},
}

type RenderTrackData struct {
	Points []gpx.GPXPoint
	Color  color.RGBA
}

func GetColorForIndex(i int) color.RGBA {
	i = i % len(ColorPalette)
	return ColorPalette[i]
}

func RenderTracks(tracks []RenderTrackData) error {

	l := log.WithFields(log.Fields{
		"comp": "gpxutils/RenderTracks",
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

func RenderHtml(tracks []RenderTrackData, filePath string, mapHtml []byte, renderPoints bool) error {

	l := log.WithFields(log.Fields{
		"comp": "gpxutils/RenderHtml",
	})

	l.Infof("Rendering to html '%s', tracks-count: %d", filePath, len(tracks))

	geoJson, err := RenderTracksDataToGeoJson(tracks, renderPoints)
	if err != nil {
		return err
	}

	tpl, err := template.New("map-html").Parse(string(mapHtml))
	if err != nil {
		return err
	}

	fhtml, err := os.Create(filePath)
	if err != nil {
		return err
	}

	err = tpl.Execute(fhtml, strings.ReplaceAll(string(geoJson), "\"", "\\\""))
	if err != nil {
		return err
	}

	fhtml.Close()

	return nil
}
