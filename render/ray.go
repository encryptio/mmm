package render

import (
	"git.encryptio.com/mmm/world"
)

type Ray struct {
	SX, SY, SZ int
	DX, DY, DZ int // over 256
}

func RayTrace(w *world.World, r Ray, includeAir, includeStart bool, maxLength int, f func(uint16, int, int, int) bool) bool {
	x := r.SX
	y := r.SY
	z := r.SZ

	dx := 0
	dy := 0
	dz := 0

	d := 0

	ld := 0
	if includeStart {
		ld = -1
	}

	for d < maxLength {
		if d != ld {
			ld = d
			id := w.Get(x, y, z)
			if includeAir || id != 0 {
				if !f(id, x, y, z) {
					return false
				}
			}
		}

		dx += r.DX
		dy += r.DY
		dz += r.DZ

		if dx >= 256 {
			dx -= 256
			x++
			d++
		}
		if dy >= 256 {
			dy -= 256
			y++
			d++
		}
		if dz >= 256 {
			dz -= 256
			z++
			d++
		}
		if dx <= -256 {
			dx += 256
			x--
			d++
		}
		if dy <= -256 {
			dy += 256
			y--
			d++
		}
		if dz <= -256 {
			dz += 256
			z--
			d++
		}

		if y >= 256 || y < 0 {
			break
		}
	}

	return true
}
