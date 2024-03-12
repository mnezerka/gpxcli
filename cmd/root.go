package cmd

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gpx",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	log.SetHandler(text.New(os.Stderr))
}

func initConfig() {}
