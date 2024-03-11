package cmd

import (
	"image/color"
	"mnezerka/gpxcli/gpxutils"

	"github.com/apex/log"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
)

var alignCmd = &cobra.Command{
	Use:   "align",
	Short: "Align two tracks",
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

		t1, t2, err := gpxutils.AlignTracks(gpx1, gpx2, -50)
		if err != nil {
			return err
		}

		// render to map
		ctx := sm.NewContext()
		ctx.SetSize(2000, 1600)

		colorOrange := color.RGBA{0x33, 0xb2, 0xff, 0xff}
		colorBlue := color.RGBA{0, 0, 0xff, 0xff}
		colorRed := color.RGBA{0xff, 0, 0, 0xff}
		colorSome := color.RGBA{0xff, 0x99, 0, 0xff}

		for i := 0; i < len(t1); i++ {
			if t1[i] != nil && t2[i] != nil {
				ctx.AddObject(sm.NewCircle(s2.LatLngFromDegrees(t1[i].Latitude, t1[i].Longitude), colorOrange, colorOrange, 10, 0))
				ctx.AddObject(sm.NewCircle(s2.LatLngFromDegrees(t2[i].Latitude, t2[i].Longitude), colorBlue, colorBlue, 10, 0))
			} else if t1[i] != nil && t2[i] == nil {
				ctx.AddObject(sm.NewCircle(s2.LatLngFromDegrees(t1[i].Latitude, t1[i].Longitude), colorRed, colorRed, 10, 0))
			} else if t1[i] == nil && t2[i] != nil {
				ctx.AddObject(sm.NewCircle(s2.LatLngFromDegrees(t2[i].Latitude, t2[i].Longitude), colorSome, colorSome, 10, 0))
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
	},
}

func init() {
	rootCmd.AddCommand(alignCmd)
}
