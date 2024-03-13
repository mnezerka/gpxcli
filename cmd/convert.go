package cmd

import (
	"mnezerka/gpxcli/gpxutils"
	"os"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
)

var (
	ConfigConvert struct {
		Output                 string
		OutputHtml             bool
		OutputGpx              bool
		Interpolate            bool
		InterpolateMaxDistance int
		Simplify               bool
		SimplifyMinDistance    int
		Reverse                bool
		ShowPoints             bool
	}
)

var convertCmd = &cobra.Command{
	Use:   "convert gpx_file_path",
	Short: "Convert gpx file - apply sequence of optional operations",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		l := log.WithFields(log.Fields{
			"comp": "cmd/convert",
		})

		l.Infof("Convert args %v", args)

		gpxFile, err := gpx.ParseFile(args[0])
		if err != nil {
			log.WithError(err)
			return err
		}

		l.Infof("Track (%s)", gpxFile.Name)
		l.Infof("  length 2D: %.1fkm", gpxFile.Length2D()/1000)
		l.Infof("  length 3D: %.1fkm", gpxFile.Length3D()/1000)
		l.Infof("  tracks: %d", len(gpxFile.Tracks))

		///////////////////////////////// simplification
		if ConfigConvert.Simplify {
			l.Infof("Simplifying, min-distance: %d", ConfigConvert.SimplifyMinDistance)
			l.Infof("  gpx-points: %d", gpxFile.GetTrackPointsNo())
			gpxFile.SimplifyTracks(float64(ConfigConvert.SimplifyMinDistance))
			l.Infof("  gpx-simplified-points: %d", gpxFile.GetTrackPointsNo())
		}

		///////////////////////////////// gpx -> list of points
		points, err := gpxutils.GpxFileToPoints(gpxFile)
		if err != nil {
			return err
		}

		///////////////////////////////// simplification
		if ConfigConvert.Reverse {
			points, err = gpxutils.Reverse(points)
			if err != nil {
				return err
			}
		}

		///////////////////////////////// distance interpolation
		if ConfigConvert.Interpolate {
			points, err = gpxutils.InterpolateDistance(points, float64(ConfigConvert.InterpolateMaxDistance))
			if err != nil {
				return err
			}
		}

		///////////////////////////////// output
		if ConfigConvert.OutputGpx {

			gpxFile.Tracks = nil

			for i := 0; i < len(points); i++ {
				gpxFile.AppendPoint(&points[i])
			}

			xmlData, err := gpxFile.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
			if err != nil {
				return err
			}

			filePath := ConfigRender.Output + ".gpx"
			err = os.WriteFile(filePath, xmlData, 0644)
			if err != nil {
				return err
			}

			l.Infof("Data written to %s\n", filePath)

		} else if ConfigRender.OutputHtml {
			mapHtml, err := GetTempalteContent("cmd/templates/map.html")
			if err != nil {
				return err
			}

			var toRender []gpxutils.RenderTrackData
			toRender = append(toRender, gpxutils.RenderTrackData{
				Points: points,
				Color:  gpxutils.GetColorForIndex(0),
			})

			err = gpxutils.RenderHtml(toRender, ConfigRender.Output+".html", mapHtml, ConfigConvert.ShowPoints)
			if err != nil {
				return err
			}
		} else {
			l.Warn("No output specified")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&ConfigConvert.Output, "output", "o", "aligned", "Output file path without extension (will be added based on output format)")
	convertCmd.Flags().BoolVar(&ConfigConvert.OutputGpx, "gpx", false, "Output in gpx format")
	convertCmd.Flags().BoolVar(&ConfigConvert.OutputHtml, "html", true, "Output in htmlformat")
	convertCmd.MarkFlagsMutuallyExclusive("gpx", "html")

	convertCmd.Flags().BoolVarP(&ConfigConvert.Interpolate, "interpolate", "i", false, "Interpolate distances")
	convertCmd.Flags().IntVarP(&ConfigConvert.InterpolateMaxDistance, "interpolate-distance", "", 10, "Maximal distance (in meters) for points interpolation")

	convertCmd.Flags().BoolVar(&ConfigConvert.Simplify, "simplify", false, "Simplify track")
	convertCmd.Flags().IntVarP(&ConfigConvert.SimplifyMinDistance, "simplify-min-distance", "", 10, "Minimal distance for simplification")

	convertCmd.Flags().BoolVar(&ConfigConvert.Reverse, "reverse", false, "Reverse track")

	convertCmd.Flags().BoolVar(&ConfigConvert.ShowPoints, "render-points", false, "Draw symbol (circle) for each track point if rendering to map")

}
