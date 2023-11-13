package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"grib-tiler/pkg/app"
	"grib-tiler/pkg/opts"
)

func main() {
	log := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
	opts := opts.Opts{}
	a := &cli.App{
		Name:      "grib-tiler",
		Usage:     "CLI app to create tiles and vector contours from raster GRIB2 files.",
		UsageText: "grib-tiler [options] input output",
		Action: func(cCtx *cli.Context) error {
			in := cCtx.Args().Get(0)
			out := cCtx.Args().Get(1)
			if _, err := os.Stat(in); err != nil {
				return fmt.Errorf("Файл %q не найден.", in)
			}
			opts.InputFile = in
			if out == "" {
				opts.OutputDirectory = "./result"
			} else {
				opts.OutputDirectory = out
			}
			return nil
		},
		Flags:   opts.Flags(),
		Version: "1.1.0",
		Authors: []*cli.Author{
			{
				Name:  "Ivan Konev",
				Email: "ivankonewv@gmail.com",
			},
		},
	}
	if err := a.Run(os.Args); err != nil {
		log.Fatal(err)
		return
	}
	if opts.Zoom == nil {
		opts.Zoom = []int{0, 1, 2}
	}

	app.NewApp(&opts, log).Run()
}
