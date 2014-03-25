package camera

import (
	"git.encryptio.com/mmm/render"
)

type Camera interface {
	RayAt(x, y int) render.Ray
}
