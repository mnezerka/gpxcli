package cmd

import (
	"fmt"
	"github.com/mnezerka/gpxcli/gpxutils"
	"image/color"
	"os"
	"strings"
	"text/template"

	"github.com/apex/log"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
)

var (
	ConfigAlign struct {
		Output              string
		OutputPng           bool
		OutputHtml          bool
		OutputYaml          bool
		MinDistance         int
		Simplify            bool
		SimplifyMinDistance int
	}
)

var colorLightBlue color.RGBA = color.RGBA{0x33, 0xb2, 0xff, 0xff}
var colorBlue color.RGBA = color.RGBA{0, 0, 0xff, 0xff}
var colorRed color.RGBA = color.RGBA{0xff, 0, 0, 0xff}
var colorOrange color.RGBA = color.RGBA{0xff, 0x99, 0, 0xff}

var alignCmd = &cobra.Command{
	Use:   "align gpx1 gpx2",
	Short: "Align two gpx tracks",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		l := log.WithFields(log.Fields{
			"comp": "cmd/align",
		})

		l.Infof("Aligning gpx files %v", args)

		gpx1, err := gpx.ParseFile(args[0])
		if err != nil {
			log.WithError(err)
			return err
		}
		gpx2, err := gpx.ParseFile(args[1])
		if err != nil {
			log.WithError(err)
			return err
		}

		if ConfigAlign.Simplify {
			l.Infof("Simplifying gpx files, min-distance: %d", ConfigAlign.SimplifyMinDistance)
			l.Infof("  gpx1-points: %d", gpx1.GetTrackPointsNo())
			l.Infof("  gpx2-points: %d", gpx2.GetTrackPointsNo())
			gpx1.SimplifyTracks(float64(ConfigAlign.SimplifyMinDistance))
			gpx2.SimplifyTracks(float64(ConfigAlign.SimplifyMinDistance))
			l.Infof("  gpx1-simplified-points: %d", gpx1.GetTrackPointsNo())
			l.Infof("  gpx2-simplified-points: %d", gpx2.GetTrackPointsNo())
		}

		t1, t2, err := gpxutils.AlignTracks(gpx1, gpx2, float64(-1*ConfigAlign.MinDistance))
		if err != nil {
			return err
		}

		if ConfigAlign.OutputPng {
			err := renderAlignedImage(t1, t2)
			if err != nil {
				return err
			}
		} else if ConfigAlign.OutputHtml {
			err := renderAlignedHtml(t1, t2)
			if err != nil {
				return err
			}
		} else if ConfigAlign.OutputYaml {
			err := renderAlignedYaml(t1, t2)
			if err != nil {
				return err
			}
		} else {
			l.Warn("No output format specified")
		}

		return nil
	},
}

func renderAlignedHtml(t1, t2 []*gpx.GPXPoint) error {

	l := log.WithFields(log.Fields{
		"comp": "cmd/align/renderAlignedHtml",
	})

	l.Infof("Rendering to html, t1 and t2 points: %d", len(t1))

	fc := geojson.NewFeatureCollection()

	var f *geojson.Feature

	l.Infof("  converting to GeoJson format")

	// go through points
	for i := 0; i < len(t1); i++ {
		// if points were matched
		if t1[i] != nil && t2[i] != nil {

			f = geojson.NewFeature(orb.Point{t1[i].Longitude, t1[i].Latitude})
			f.Properties["color"] = fmt.Sprintf("#%02x%02x%02x", colorLightBlue.R, colorLightBlue.G, colorLightBlue.B)
			fc.Append(f)

			f = geojson.NewFeature(orb.Point{t2[i].Longitude, t2[i].Latitude})
			f.Properties["color"] = fmt.Sprintf("#%02x%02x%02x", colorBlue.R, colorBlue.G, colorBlue.B)
			fc.Append(f)

			// connection line between matched points
			f = geojson.NewFeature(orb.LineString{
				orb.Point{t1[i].Longitude, t1[i].Latitude},
				orb.Point{t2[i].Longitude, t2[i].Latitude},
			})
			f.Properties["color"] = fmt.Sprintf("#%02x%02x%02x", colorBlue.R, colorBlue.G, colorBlue.B)
			fc.Append(f)

		} else if t1[i] != nil && t2[i] == nil {
			f = geojson.NewFeature(orb.Point{t1[i].Longitude, t1[i].Latitude})
			f.Properties["color"] = fmt.Sprintf("#%02x%02x%02x", colorRed.R, colorRed.G, colorRed.B)
			fc.Append(f)
		} else if t1[i] == nil && t2[i] != nil {

			f = geojson.NewFeature(orb.Point{t2[i].Longitude, t2[i].Latitude})
			f.Properties["color"] = fmt.Sprintf("#%02x%02x%02x", colorOrange.R, colorOrange.G, colorOrange.B)
			fc.Append(f)
		}
	}

	geoJson, err := fc.MarshalJSON()
	if err != nil {
		return err
	}

	tplFilePath := "cmd/templates/map.html"

	l.Infof("  reading html template '%s'", tplFilePath)

	mapHtml, err := templatesContent.ReadFile(tplFilePath)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return fmt.Errorf("template '%s' does not exist", "map.html")
		}

		// unknown error
		return err
	}

	tpl, err := template.New("map-html").Parse(string(mapHtml))
	if err != nil {
		return err
	}

	filePath := ConfigRender.Output + ".html"

	l.Infof("  writing %d bytes to '%s'", len(geoJson), filePath)

	fhtml, err := os.Create(filePath)
	if err != nil {
		return err
	}

	err = tpl.Execute(fhtml, strings.ReplaceAll(string(geoJson), "\"", "\\\""))
	if err != nil {
		return err
	}

	fhtml.Close()

	l.Info("  done")

	return nil
}

func renderAlignedImage(t1, t2 []*gpx.GPXPoint) error {

	// render to map
	ctx := sm.NewContext()
	ctx.SetSize(2000, 1600)

	for i := 0; i < len(t1); i++ {
		if t1[i] != nil && t2[i] != nil {
			ctx.AddObject(sm.NewCircle(s2.LatLngFromDegrees(t1[i].Latitude, t1[i].Longitude), colorLightBlue, colorLightBlue, 10, 0))
			ctx.AddObject(sm.NewCircle(s2.LatLngFromDegrees(t2[i].Latitude, t2[i].Longitude), colorBlue, colorBlue, 10, 0))
		} else if t1[i] != nil && t2[i] == nil {
			ctx.AddObject(sm.NewCircle(s2.LatLngFromDegrees(t1[i].Latitude, t1[i].Longitude), colorRed, colorRed, 10, 0))
		} else if t1[i] == nil && t2[i] != nil {
			ctx.AddObject(sm.NewCircle(s2.LatLngFromDegrees(t2[i].Latitude, t2[i].Longitude), colorOrange, colorOrange, 10, 0))
		}
	}

	img, err := ctx.Render()
	if err != nil {
		return err
	}

	if err := gg.SavePNG(ConfigAlign.Output+".png", img); err != nil {
		return err
	}

	return nil
}

func renderAlignedYaml(t1, t2 []*gpx.GPXPoint) error {

	var countUnique1, countUnique2, countMatching int

	for i := 0; i < len(t1); i++ {
		if t1[i] != nil && t2[i] != nil {
			countMatching++
		} else if t1[i] != nil && t2[i] == nil {
			countUnique1++
		} else if t1[i] == nil && t2[i] != nil {
			countUnique2++
		}
	}

	var matchGuess float64 = (float64(countMatching) * 100.0) / float64(len(t1))

	fmt.Printf("total-points: %d\n", len(t1))
	fmt.Printf("matching-points: %d\n", countMatching)
	fmt.Printf("track1-unique-points: %d\n", countUnique1)
	fmt.Printf("track2-unique-points: %d\n", countUnique2)
	fmt.Printf("match-guess: %.1f%%\n", matchGuess)
	return nil
}

func init() {
	rootCmd.AddCommand(alignCmd)

	alignCmd.Flags().StringVarP(&ConfigAlign.Output, "output", "o", "aligned", "Output file path without extension (will be added based on output format)")

	alignCmd.Flags().BoolVar(&ConfigAlign.OutputPng, "png", false, "Output to png image")
	alignCmd.Flags().BoolVar(&ConfigAlign.OutputHtml, "html", false, "Output to html page")
	alignCmd.Flags().BoolVar(&ConfigAlign.OutputYaml, "yaml", true, "Output to yaml format")
	alignCmd.MarkFlagsMutuallyExclusive("png", "html", "yaml")

	alignCmd.Flags().IntVarP(&ConfigAlign.MinDistance, "min-distance", "", 50, "Minimal distance between points in meters")

	alignCmd.Flags().BoolVar(&ConfigAlign.Simplify, "simplify", false, "Simplify gpx tracks before aligning")
	alignCmd.Flags().IntVarP(&ConfigAlign.SimplifyMinDistance, "simplify-min-distance", "", 10, "Minimal distance for simplification")
}
