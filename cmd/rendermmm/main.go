package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"git.encryptio.com/mmm/camera"
	"git.encryptio.com/mmm/render"
	"git.encryptio.com/mmm/world"
)

func main() {
	sx := flag.Float64("x", -63, "Starting x position")
	sz := flag.Float64("z", -63, "Starting z position")
	wid := flag.Int("w", 128, "Width of output image")
	hei := flag.Int("h", 128, "Height of output image")
	regionPath := flag.String("region", "hub/region/", "Path to region directory of world")
	outputPath := flag.String("o", "out.png", "Output filename")
	day := flag.Bool("day", true, "Daytime")
	lightSamples := flag.Int("lightsamples", 100, "Number of light samples per probe")
	pixelSamples := flag.Int("pixelsamples", 1, "Square root of number of probes per pixel")
	cameraType := flag.String("camera", "iso", "Camera type (iso, topdown)")
	emptyExit := flag.Int("emptyexit", 0, "Exit code to return when the image is completely empty")
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

	sampleStep := 1.0 / float64(*pixelSamples)

	im := image.NewNRGBA(image.Rect(0, 0, *wid, *hei))
	for x := 0; x < *wid; x++ {
		for y := 0; y < *hei; y++ {
			var co color.NRGBA64
			for dx := float64(0); dx < 0.999; dx += sampleStep {
				for dy := float64(0); dy < 0.999; dy += sampleStep {
					part := render.ShadeRay(w, c.RayAt(float64(x)+dx, float64(y)+dy), sky, *lightSamples)
					co.R += uint16(part.R)
					co.G += uint16(part.G)
					co.B += uint16(part.B)
					co.A += uint16(part.A)
				}
			}

			co.R /= uint16(*pixelSamples * *pixelSamples)
			co.G /= uint16(*pixelSamples * *pixelSamples)
			co.B /= uint16(*pixelSamples * *pixelSamples)
			co.A /= uint16(*pixelSamples * *pixelSamples)

			im.SetNRGBA(x, y, color.NRGBA{uint8(co.R), uint8(co.G), uint8(co.B), uint8(co.A)})
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

	isEmpty := true
	for x := 0; x < *wid; x++ {
		for y := 0; y < *hei; y++ {
			if im.Pix[im.PixOffset(x, y)+3] != 0 {
				isEmpty = false
				break
			}
		}

		if !isEmpty {
			break
		}
	}

	if isEmpty {
		os.Exit(*emptyExit)
	} else {
		os.Exit(0)
	}
}
