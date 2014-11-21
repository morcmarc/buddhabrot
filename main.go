package main

import (
	"log"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Buddhabrot"
	app.Usage = "generate buddahbrot fractals"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "width, w",
			Value: 600,
		},
		cli.IntFlag{
			Name:  "height, t",
			Value: 600,
		},
		cli.StringFlag{
			Name:  "color, c",
			Usage: "color palette, if only one value given (i.e.: '-c 40') it will render a greyscale image",
			Value: "50,200,500",
		},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			log.Fatalln("Missing output file")
		}
		oFileName := c.Args()[0]
		oFile, err := os.Create(oFileName)
		if err != nil {
			log.Fatalf("Could not create file %s, %s", oFileName, err)
		}
		defer oFile.Close()

		width := c.Int("width")
		height := c.Int("height")
		colors := SplitColors(c.String("color"))

		buddha := NewBuddha(width, height, colors, oFile)
		sdlHandler := NewSdlHandler(width, height, buddha)

		cpus := runtime.NumCPU()
		runtime.GOMAXPROCS(cpus)

		sdlHandler.Start()
	}
	app.Run(os.Args)
}
