package render

import (
	"git.encryptio.com/mmm/world"

	"image/color"
	"log"
	"math/rand"
)

var SkyDay = color.NRGBA{255, 255, 255, 255}
var SkyNight = color.NRGBA{32, 32, 64, 255}

func over(bottom, top color.NRGBA) (ret color.NRGBA) {
	a := uint16(top.A)
	b := (uint16(bottom.A) * (255 - a)) / 255
	ret.R = uint8((uint16(top.R)*a + uint16(bottom.R)*b) / 255)
	ret.G = uint8((uint16(top.G)*a + uint16(bottom.G)*b) / 255)
	ret.B = uint8((uint16(top.B)*a + uint16(bottom.B)*b) / 255)
	ret.A = uint8(a + b)
	return ret
}

func ShadeRay(w *world.World, ray Ray, sky color.NRGBA, samples int) color.NRGBA {
	c := color.NRGBA{0, 0, 0, 0}

	var firstX, firstY, firstZ int
	var hit bool

	RayTrace(w, ray, false, true, 1024, func(id uint16, x, y, z int) bool {
		if !hit {
			hit = true
			firstX = x
			firstY = y
			firstZ = z
		}

		blockInfo, ok := BlockInfo[id]
		if !ok {
			log.Printf("No color for block id %d\n", id)
			return true
		}

		c = over(blockInfo.Color, c)
		return c.A != 255
	})

	if hit {
		light := AmbientAt(w, firstX, firstY+1, firstZ, sky, samples)
		c.R = uint8(uint16(c.R) * uint16(light.R) / 255)
		c.G = uint8(uint16(c.G) * uint16(light.G) / 255)
		c.B = uint8(uint16(c.B) * uint16(light.B) / 255)
	}

	return c
}

func AmbientAt(w *world.World, sx, sy, sz int, sky color.NRGBA, samples int) (ret color.NRGBA) {
	ret = RadiosityAt(w, sx, sy, sz, sky, samples)
	ret.R = ret.R/2 + (ret.R/2+sky.R/2)/2
	ret.G = ret.G/2 + (ret.G/2+sky.G/2)/2
	ret.B = ret.B/2 + (ret.B/2+sky.B/2)/2
	return
}

func RadiosityAt(w *world.World, sx, sy, sz int, sky color.NRGBA, samples int) color.NRGBA {
	r := float64(0)
	g := float64(0)
	b := float64(0)

	for i := 0; i < samples; i++ {
		dx := rand.Int()%128 - 64
		dy := rand.Int() % 64
		dz := rand.Int()%128 - 64

		for dx*dx+dy*dy+dz*dz > 64*64*3 || dx*dx+dy*dy+dz*dz < 32*32*3 {
			dx = rand.Int()%128 - 64
			dy = rand.Int() % 64
			dz = rand.Int()%128 - 64
		}

		hitEnd := RayTrace(w, Ray{sx, sy, sz, dx, dy, dz}, false, false, 96, func(id uint16, x, y, z int) bool {
			info, ok := BlockInfo[id]
			if !ok {
				log.Printf("No color for block id %d\n", id)
				return true
			}

			if info.Emission > 0 {
				r += float64(info.Color.R) * float64(info.Emission) / 255 * 16
				g += float64(info.Color.G) * float64(info.Emission) / 255 * 16
				b += float64(info.Color.B) * float64(info.Emission) / 255 * 16
			}

			return false
		})

		if hitEnd {
			r += float64(sky.R)
			g += float64(sky.G)
			b += float64(sky.B)
		}
	}

	r /= float64(samples)
	g /= float64(samples)
	b /= float64(samples)

	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}

	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}
