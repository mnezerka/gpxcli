package main

import (
	"embed"
	"fmt"
	"mnezerka/gpxcli/cmd"
	"os"

	"github.com/apex/log"
)

//go:embed cmd/templates/*
var templatesContent embed.FS

func main() {

	cmd.SetTemplatesContent(&templatesContent)

	if err := cmd.Execute(); err != nil {
		log.WithError(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
