// "github.com/tkrajina/gpxgo/gpx"

package cmd

import (
	"mnezerka/gpxcli/gpxutils"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
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
			l.Infof("  - length 2D: %f", gpx.Length2D())
			l.Infof("  - length 3D: %f", gpx.Length3D())
			l.Infof("  - tracks: %d", len(gpx.Tracks))

			points, err := gpxutils.GpxFileToPoints(gpx)
			if err != nil {
				return err
			}

			toRender = append(toRender, gpxutils.RenderTrackData{
				Points: points,
				Color:  gpxutils.GetColorForIndex(file_ix),
			})

			gpxutils.RenderTracks(toRender)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
}
