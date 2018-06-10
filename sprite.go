package wssprite

import (
	"github.com/oshimaya/gowsdisplay"
	"image"
)

type Pattern struct {
	center image.Point
	pix    gowsdisplay.PIXELARRAY
}

type Sprite struct {
	id   int
	pos  image.Point
	view bool
	hit  bool
	pat  *Pattern
}

func NewSprite(id int) (sp Sprite) {
	sp.id = id
	sp.pos = image.Pt(0, 0)
	sp.view = false
	return sp
}

func (sp *Sprite) SetPosition(x int, y int) {
	sp.pos = image.Pt(x, y)
}

func (sp *Sprite) GetPos() (int, int) {
	return sp.pos.X, sp.pos.Y
}

func (sp *Sprite) Enable() {
	sp.view = true
}

func (sp *Sprite) Disable() {
	sp.view = false
}

func (sp *Sprite) SetPattern(pat *Pattern) {
	sp.pat = pat
	sp.pat.center = image.Pt(pat.pix.GetWidth()/2, pat.pix.GetHeight()/2)
}
