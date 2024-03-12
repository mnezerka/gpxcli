// "github.com/tkrajina/gpxgo/gpx"

package cmd

import (
	"fmt"
	"mnezerka/gpxcli/gpxutils"
	"os"
	"strings"
	"text/template"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
)

var (
	ConfigRender struct {
		Output     string
		OutputPng  bool
		OutputHtml bool
		Points     bool
	}
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render tracks to map",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		l := log.WithFields(log.Fields{
			"comp": "cmd/render",
		})

		var toRender []gpxutils.RenderTrackData

		for file_ix := 0; file_ix < len(args); file_ix++ {

			gpx, err := gpx.ParseFile(args[file_ix])
			if err != nil {
				log.WithError(err)
				return err
			}

			l.Infof("Track (%s)", gpx.Name)
			l.Infof("  length 2D: %.1fkm", gpx.Length2D()/1000)
			l.Infof("  length 3D: %.1fkm", gpx.Length3D()/1000)
			l.Infof("  tracks: %d", len(gpx.Tracks))

			points, err := gpxutils.GpxFileToPoints(gpx)
			if err != nil {
				return err
			}

			toRender = append(toRender, gpxutils.RenderTrackData{
				Points: points,
				Color:  gpxutils.GetColorForIndex(file_ix),
			})
		}

		if ConfigRender.OutputPng {
			err := gpxutils.RenderTracks(toRender)
			if err != nil {
				return err
			}
		} else if ConfigRender.OutputHtml {
			err := renderHtml(toRender)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func renderHtml(tracks []gpxutils.RenderTrackData) error {

	mapHtml, err := templatesContent.ReadFile("cmd/templates/map.html")
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return fmt.Errorf("template '%s' does not exist", "map.html")
		}

		// unknown error
		return err
	}

	geoJson, err := gpxutils.RenderTracksDataToGeoJson(tracks, ConfigRender.Points)
	if err != nil {
		return err
	}

	tpl, err := template.New("map-html").Parse(string(mapHtml))
	if err != nil {
		return err
	}

	fhtml, err := os.Create(ConfigRender.Output + ".html")
	if err != nil {
		return err
	}

	err = tpl.Execute(fhtml, strings.ReplaceAll(string(geoJson), "\"", "\\\""))
	if err != nil {
		return err
	}

	fhtml.Close()

	/*err = os.WriteFile(ConfigRender.Output+".html", mapHtml, 0644)
	if err != nil {
		return err
	}


	// write geojson to file
	err = os.WriteFile(ConfigRender.Output+".json", geoJson, 0644)
	if err != nil {
		return err
	}
	*/

	return nil
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVarP(&ConfigRender.Output, "output", "o", "map", "Output file path without extension (will be added based on output format)")

	renderCmd.Flags().BoolVar(&ConfigRender.OutputPng, "png", false, "Output in PNG format")
	renderCmd.Flags().BoolVar(&ConfigRender.OutputHtml, "html", true, "Output in HTML format")
	renderCmd.MarkFlagsMutuallyExclusive("png", "html")

	renderCmd.Flags().BoolVar(&ConfigRender.Points, "points", false, "Draw symbol (circle) for each track point")

}
