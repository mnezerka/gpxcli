package cmd

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
)

// params (flags)
var (
	config_simplify_min_distance int
	config_simplify_output       string
)

var simplifyCmd = &cobra.Command{
	Use:   "simplify",
	Short: "Simplify tracks in single gpx file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("simplifying file %v", args[0])

		gpxFile, err := gpx.ParseFile(args[0])
		if err != nil {
			log.WithError(err)
			return err
		}

		gpxFile.SimplifyTracks(float64(config_simplify_min_distance))

		xmlData, err := gpxFile.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
		if err != nil {
			return err
		}

		err = os.WriteFile(config_simplify_output, xmlData, 0644)
		if err != nil {
			return err
		}

		fmt.Printf("simplified data written to %s\n", config_simplify_output)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(simplifyCmd)
	simplifyCmd.Flags().IntVarP(&config_simplify_min_distance, "min-distance", "", 5, "Minimal distance in meters for track simplification")
	simplifyCmd.Flags().StringVarP(&config_simplify_output, "output", "o", "output.gpx", "Output file (including extension)")
}
