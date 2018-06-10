package wssprite

import (
	"github.com/oshimaya/gowsdisplay"
	"image"
)

type BgScreen struct {
	id           int
	row          int
	column       int
	cell_width   int
	cell_height  int
	offset       image.Point
	h_scrollable bool
	v_scrollable bool
	priority     int
	text         []int
	pat          []gowsdisplay.PIXELARRAY
}

func NewBg(id int, row int, col int, cw int, ch int) *BgScreen {
	bg := new(BgScreen)
	bg.id = id
	bg.row = row
	bg.column = col
	bg.offset = image.Pt(0, 0)
	bg.h_scrollable = true
	bg.v_scrollable = true
	bg.text = make([]int, row*col)
	bg.pat = make([]gowsdisplay.PIXELARRAY, 0)
	return bg
}

func (bg *BgScreen) SetOffset(x int, y int) {
	bg.offset = image.Pt(x, y)
}

func (bg *BgScreen) GetOffset() (int, int) {
	return bg.offset.X, bg.offset.Y
}

func (bg *BgScreen) SetPriority(p int) {
	bg.priority = p
}

func (bg *BgScreen) SetScrollable(h, v bool) {
	bg.h_scrollable = h
	bg.v_scrollable = v
}

func (bg *BgScreen) Put(x int, y int, patid int) {
	if patid < 0 || patid >= len(bg.pat) ||
		x < 0 || x >= bg.row || y < 0 || y >= bg.column {
		// Nothing to do.
		return
	}
	bg.text[x+y*bg.row] = patid
}

func (bg *BgScreen) Get(x int, y int) (patid int) {
	if x < 0 || x >= bg.row || y < 0 || y >= bg.column {
		// Nothing to do.
		return -1
	}
	return bg.text[x+y*bg.row]
}

func (bg *BgScreen) AddPixelPattern(pix gowsdisplay.PIXELARRAY) {
	if pix.GetWidth() != bg.cell_width || pix.GetHeight() != bg.cell_height {
		return
	}
	bg.pat = append(bg.pat, pix)
}
