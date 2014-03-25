package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	"git.encryptio.com/mmm/camera"
	"git.encryptio.com/mmm/render"
	"git.encryptio.com/mmm/world"
)

func main() {
	sx := flag.Int("x", -63, "Starting x position")
	sz := flag.Int("z", -63, "Starting z position")
	wid := flag.Int("w", 128, "Width of output image")
	hei := flag.Int("h", 128, "Height of output image")
	regionPath := flag.String("region", "hub/region/", "Path to region directory of world")
	outputPath := flag.String("o", "out.png", "Output filename")
	flag.Parse()

	w := world.New(*regionPath)

	c := camera.TopDown(*sx, *sz)

	im := image.NewNRGBA(image.Rect(0, 0, *wid, *hei))
	for x := 0; x < *wid; x++ {
		log.Printf("%d/%d\n", x, *wid)
		for y := 0; y < *hei; y++ {
			im.SetNRGBA(x, y, render.ShadeRay(w, c.RayAt(x, y), render.SkyNight))
		}
	}

	fh, err := os.Create(*outputPath)
	if err != nil {
		log.Printf("Couldn't open %s for writing: %s\n", *outputPath, err.Error())
		os.Exit(1)
	}

	err = png.Encode(fh, im)
	if err != nil {
		log.Printf("Couldn't encode to PNG: %s\n", err.Error())
		os.Exit(1)
	}

	err = fh.Close()
	if err != nil {
		log.Printf("Couldn't close filehandle: %s\n", err.Error())
		os.Exit(1)
	}
}
