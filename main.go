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
	day := flag.Bool("day", true, "Daytime")
	samples := flag.Int("samples", 100, "Number of samples per pixel")
	cameraType := flag.String("camera", "iso", "Camera type (iso, topdown)")
	flag.Parse()

	w := world.New(*regionPath)
	sky := render.SkyDay
	if !*day {
		sky = render.SkyNight
	}

	var c camera.Camera
	switch *cameraType {
	case "iso":
		c = camera.Isometric(*sx, *sz)
	case "topdown":
		c = camera.TopDown(*sx, *sz)
	default:
		log.Printf("Unknown camera type %v", *cameraType)
		os.Exit(1)
	}

	im := image.NewNRGBA(image.Rect(0, 0, *wid, *hei))
	for x := 0; x < *wid; x++ {
		log.Printf("%d/%d\n", x, *wid)
		for y := 0; y < *hei; y++ {
			im.SetNRGBA(x, y, render.ShadeRay(w, c.RayAt(x, y), sky, *samples))
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
