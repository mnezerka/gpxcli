# GpxCli

![Release Build](https://github.com/mnezerka/gpxcli/actions/workflows/release.yml/badge.svg?event=release)

Tool for mainipulation with gpx files

## Merge

Allows to merge content from multiple gpx files into one single file. Merge
operation does also simplification of track data

*Use case - generating heat maps from hundreads of tracks in your
browser is more efficient if you download preprocessed data in one file instead of
fetching each gpx file separately and parsing from xml.*

Examples:

Merge three gpx tracks into single file `output.json` (default output format is json)
```bash
./gpxcli merge track1.gpx track2.gpx track3.gpx
```

Merge all gpx files from gpx directory to file `all.yml` with  custom value for
minimal distance (in megetrs) between points (track simplification):
```bash
./gpxcli merge --min-distance 20 --output all --yaml gpx/*
```
