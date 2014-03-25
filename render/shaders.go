package render

import (
	"git.encryptio.com/mmm/world"

	"image/color"
	"log"
	"math/rand"
)

var SkyDay = color.NRGBA{255, 255, 255, 255}
var SkyNight = color.NRGBA{24, 24, 48, 255}

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
	var accum color.NRGBA64

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
				accum.R += uint16(info.Color.R) * uint16(info.Emission) / 255 * 16
				accum.G += uint16(info.Color.G) * uint16(info.Emission) / 255 * 16
				accum.B += uint16(info.Color.B) * uint16(info.Emission) / 255 * 16
			}

			return false
		})

		if hitEnd {
			accum.R += uint16(sky.R)
			accum.G += uint16(sky.G)
			accum.B += uint16(sky.B)
		}
	}

	return color.NRGBA{uint8(accum.R / uint16(samples)), uint8(accum.G / uint16(samples)), uint8(accum.B / uint16(samples)), 255}
}
