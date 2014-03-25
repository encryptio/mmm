package camera

import (
	"git.encryptio.com/mmm/render"
)

func TopDown(offX, offZ int) Camera {
	return Ortho(offX, offZ, 1, 0, 0, 1, 0, -256, 0)
}

func Isometric(offX, offZ int) Camera {
	return Ortho(offX, offZ, 1, -1, 1, 1, 32, -64, -32)
}

func Ortho(offX, offZ int, xdirX, xdirZ, zdirX, zdirZ float64, dirX, dirY, dirZ int) Camera {
	return &ortho{offX, offZ, xdirX, xdirZ, zdirX, zdirZ, dirX, dirY, dirZ}
}

type ortho struct {
	offX, offZ       int
	xdirX, xdirZ     float64
	zdirX, zdirZ     float64
	dirX, dirY, dirZ int
}

func (c *ortho) RayAt(x, y int) render.Ray {
	x += c.offX
	z := y + c.offZ
	return render.Ray{
		SX: int(c.xdirX*float64(x) + c.xdirZ*float64(z)),
		SY: 255,
		SZ: int(c.zdirX*float64(x) + c.zdirZ*float64(z)),
		DX: c.dirX,
		DY: c.dirY,
		DZ: c.dirZ,
	}
}
