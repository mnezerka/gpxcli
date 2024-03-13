// "github.com/tkrajina/gpxgo/gpx"

package cmd

import (
	"mnezerka/gpxcli/gpxutils"

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
			mapHtml, err := GetTempalteContent("cmd/templates/map.html")
			if err != nil {
				return err
			}

			err = gpxutils.RenderHtml(toRender, ConfigRender.Output+".html", mapHtml, ConfigRender.Points)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVarP(&ConfigRender.Output, "output", "o", "map", "Output file path without extension (will be added based on output format)")

	renderCmd.Flags().BoolVar(&ConfigRender.OutputPng, "png", false, "Output in PNG format")
	renderCmd.Flags().BoolVar(&ConfigRender.OutputHtml, "html", true, "Output in HTML format")
	renderCmd.MarkFlagsMutuallyExclusive("png", "html")

	renderCmd.Flags().BoolVar(&ConfigRender.Points, "points", false, "Draw symbol (circle) for each track point")

}
