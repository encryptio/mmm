package camera

import (
	"git.encryptio.com/mmm/render"
)

type Camera interface {
	RayAt(x, y float64) render.Ray
}
