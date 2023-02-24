// "github.com/tkrajina/gpxgo/gpx"

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
	"gopkg.in/yaml.v3"
)

type SimpleTracks struct {
	Tracks []SimpleTrack `json:"tracks"`
}

type SimpleTrack struct {
	Points []SimplePoint `json:"points"`
}

type SimplePoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// params (flags)
var (
	config_merge_min_distance int
	config_merge_output       string
	config_merge_output_json  bool
	config_merge_output_yaml  bool
	config_merge_output_gpx   bool
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge tracks from multiple gpx files into one",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Do Stuff Here
		log.Infof("merging files %v", args)

		all := SimpleTracks{}

		count_points_total_raw := 0
		count_points_total_simplified := 0

		for file_ix := 0; file_ix < len(args); file_ix++ {

			fmt.Println(args[file_ix])

			gpx, err := gpx.ParseFile(args[file_ix])
			if err != nil {
				log.WithError(err)
				return err
			}
			fmt.Printf("  - length 2D: %f\n", gpx.Length2D())
			fmt.Printf("  - length 3D: %f\n", gpx.Length3D())
			fmt.Printf("  - tracks: %d\n", len(gpx.Tracks))

			gpx.ReduceGpxToSingleTrack()

			count_points_raw := 0
			for _, segment := range gpx.Tracks[0].Segments {
				count_points_raw += len(segment.Points)
			}
			count_points_total_raw += count_points_raw

			gpx.SimplifyTracks(float64(config_merge_min_distance))

			count_points_simplified := 0
			for _, segment := range gpx.Tracks[0].Segments {
				count_points_simplified += len(segment.Points)
			}
			count_points_total_simplified += count_points_simplified

			fmt.Printf("  - points simplification: %d -> %d\n", count_points_raw, count_points_simplified)

			for _, track := range gpx.Tracks {
				t := SimpleTrack{}
				for _, segment := range track.Segments {
					for _, point := range segment.Points {
						p := SimplePoint{}
						p.Lat = point.Latitude
						p.Lng = point.Longitude
						t.Points = append(t.Points, p)
					}
				}
				all.Tracks = append(all.Tracks, t)
			}
		}

		var data []byte
		var err error
		output := config_merge_output
		if config_merge_output_yaml {
			data, err = yaml.Marshal(&all)
			if err != nil {
				return err
			}
			output += ".yml"
		} else if config_merge_output_json {
			data, err = json.MarshalIndent(&all, "", "  ")
			if err != nil {
				return err
			}
			output += ".json"
		} else if config_merge_output_gpx {
			return errors.New("gpx format output not implemented yet")
		}
		err = ioutil.WriteFile(output, data, 0644)
		if err != nil {
			return err
		}

		fmt.Printf("merged data written to %s\n", output)
		fmt.Printf("  - points processed: %d\n", count_points_total_raw)
		fmt.Printf("  - points stored (after simplification): %d\n", count_points_total_simplified)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().IntVarP(&config_merge_min_distance, "min-distance", "", 50, "Minimal distance in meters for track simplification")
	mergeCmd.Flags().StringVarP(&config_merge_output, "output", "o", "output", "Output file path without extension (will be added based on output format)")

	mergeCmd.Flags().BoolVar(&config_merge_output_json, "json", true, "Output in JSON format")
	mergeCmd.Flags().BoolVar(&config_merge_output_yaml, "yaml", false, "Output in YAML format")
	mergeCmd.Flags().BoolVar(&config_merge_output_gpx, "gpx", false, "Output in GPX format")
	mergeCmd.MarkFlagsMutuallyExclusive("json", "yaml", "gpx")
}
