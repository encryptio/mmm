package camera

import (
	"git.encryptio.com/mmm/render"
)

func TopDown(offX, offZ float64) Camera {
	return Ortho(offX, offZ, 1, 0, 0, 1, 0, -256, 0)
}

func Isometric(offX, offZ float64) Camera {
	return Ortho(offX, offZ, 1, -2, 1, 2, 32, -32, -32)
}

func Ortho(offX, offZ float64, xdirX, xdirZ, zdirX, zdirZ float64, dirX, dirY, dirZ int) Camera {
	return &ortho{offX, offZ, xdirX, xdirZ, zdirX, zdirZ, dirX, dirY, dirZ}
}

type ortho struct {
	offX, offZ       float64
	xdirX, xdirZ     float64
	zdirX, zdirZ     float64
	dirX, dirY, dirZ int
}

func (c *ortho) RayAt(x, y float64) render.Ray {
	x += c.offX
	z := y + c.offZ
	return render.Ray{
		SX: int(c.xdirX*x + c.xdirZ*z),
		SY: 255,
		SZ: int(c.zdirX*x + c.zdirZ*z),
		DX: c.dirX,
		DY: c.dirY,
		DZ: c.dirZ,
	}
}
