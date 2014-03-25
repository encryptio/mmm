package camera

import (
	"git.encryptio.com/mmm/render"
)

func TopDown(offX, offZ int) Camera {
	return Ortho(offX, offZ, 1, 0, 0, 1, 0, -256, 0)
}

func Ortho(offX, offZ int, xdirX, xdirZ, zdirX, zdirZ int, dirX, dirY, dirZ int) Camera {
	return &ortho{offX, offZ, xdirX, xdirZ, zdirX, zdirZ, dirX, dirY, dirZ}
}

type ortho struct {
	offX, offZ       int
	xdirX, xdirZ     int
	zdirX, zdirZ     int
	dirX, dirY, dirZ int
}

func (c *ortho) RayAt(x, y int) render.Ray {
	x += c.offX
	z := y + c.offZ
	return render.Ray{
		SX: c.xdirX*x + c.zdirX*x,
		SY: 255,
		SZ: c.zdirX*z + c.zdirZ*z,
		DX: c.dirX,
		DY: c.dirY,
		DZ: c.dirZ,
	}
}
