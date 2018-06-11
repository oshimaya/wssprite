package wssprite

import (
	"github.com/oshimaya/gowsdisplay"
	"image"
	"image/color"
)

// Sprite screen manager
type SpriteScreen struct {
	wsd  *gowsdisplay.WsDisplay // Real Display Manage
	vsc  *Vscreen               // Virtual Screen
	sp   []Sprite               // Sprite slice
	bg   []BgScreen             // Back ground slice
	back gowsdisplay.PIXELARRAY // Lowest back ground
}

// Create new sprite screen
//	dev : Display device
//	width: Virtual screen width
//	height: Virtual screen height
//	spnum: Number of sprite
//	bgnum: Number of Background screen
//
func NewSpriteScreen(dev string, width int, height int, spnum int, bgnum int) (*SpriteScreen, error) {
	scr := new(SpriteScreen)
	wsd := gowsdisplay.NewWsDisplay(dev)
	err := wsd.Open()
	if err != nil {
		return nil, err
	}
	scr.wsd = wsd
	err = wsd.InitGraphics()
	if err != nil {
		return nil, err
	}
	scr.vsc = NewVscreen(width, height)
	pix, err := wsd.NewPixelArray()
	if err != nil {
		wsd.Close()
		return nil, err
	}
	pix.StoreImage(image.NewRGBA(image.Rect(0, 0, width, height)), wsd.GetRGBmask())

	lowpix, err := wsd.NewPixelArray()
	if err != nil {
		wsd.Close()
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < width; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{0, 0, 0, 255})
		}
	}
	lowpix.StoreImage(img, wsd.GetRGBmask())
	scr.back = lowpix
	scr.vsc.data = pix
	sp := make([]Sprite, spnum)
	for i := range sp {
		sp[i] = NewSprite(i + 1)
	}
	scr.sp = sp
	scr.bg = make([]BgScreen, 0, bgnum)
	// CreateBg

	return scr, nil
}

func (scr *SpriteScreen) Close() {
	scr.wsd.Close()
}

func (scr *SpriteScreen) GetSprite(num int) *Sprite {
	return &scr.sp[num-1]
}

func (scr *SpriteScreen) GetSprites() []Sprite {
	return scr.sp
}

func (scr *SpriteScreen) GetBg(num int) *BgScreen {
	return &scr.bg[num-1]
}

func (scr *SpriteScreen) GetVscreen() *Vscreen {
	return scr.vsc
}

func (scr *SpriteScreen) CreateBg(row int, col int, cw int, ch int) {
	id := len(scr.bg)
	bg := NewBg(id, row, col, cw, ch)

	pix, err := scr.wsd.NewPixelArray()
	if err != nil {
		return
	}
	pix.StoreImage(image.NewRGBA(image.Rect(0, 0, cw, ch)), scr.wsd.GetRGBmask())
	bg.AddPixelPattern(pix)
	scr.bg = append(scr.bg, *bg)
}

func (scr *SpriteScreen) CreatePattern(w int, h int) *Pattern {
	pat := new(Pattern)
	pix, err := scr.wsd.NewPixelArray()
	if err != nil {
		return nil
	}
	pat.pix = pix
	return pat
}

func (scr *SpriteScreen) CreatePatternFromImage(img image.Image) *Pattern {
	pat := new(Pattern)
	pix, err := scr.wsd.NewPixelArray()
	if err != nil {
		return nil
	}
	pix.StoreImage(img, scr.wsd.GetRGBmask())
	pat.pix = pix
	return pat
}

type Vscreen struct {
	offset image.Point
	data   gowsdisplay.PIXELARRAY
	attr   AttributeScreen
}

func (scr *SpriteScreen) DrawVscreen() {
	scr.wsd.PutPixelArray(scr.vsc.offset.X, scr.vsc.offset.Y, scr.vsc.data)
}

func (scr *SpriteScreen) DrawAllSprite() {
	for id := range scr.sp {
		scr.DrawSprite(id + 1)
	}
}

func (scr *SpriteScreen) DrawSprite(id int) {
	sp := scr.sp[id-1]
	if sp.view {
		sp_x := sp.pos.X - sp.pat.center.X
		sp_y := sp.pos.Y - sp.pat.center.Y
		scr.vsc.data.PutPixelPat(sp_x, sp_y, sp.pat.pix)
		if sp.hit {
			scr.vsc.SetAttr(sp_x, sp_y, sp.id)
		}
	}
}
func (scr *SpriteScreen) DrawBg(id int) {
	if id > len(scr.bg) || id <= 0 {
		return
	}
	bg := scr.bg[id-1]
	if !bg.scrollable {
		for cy := 0; cy < bg.column; cy++ {
			for cx := 0; cx < bg.row; cx++ {
				px := cx*bg.cell_width + bg.offset.X
				py := cy*bg.cell_height + bg.offset.Y
				num := bg.data[cx+cy*bg.row] // pat num
				if num < len(bg.pat) || num != 0 {
					pat := bg.pat[num]
					scr.vsc.data.PutPixelPat(px, py, pat)
				}
			}
		}
	} else {
		dx, ex := caldiff(bg.row, bg.cell_width, bg.offset.X)
		dy, ey := caldiff(bg.column, bg.cell_height, bg.offset.Y)
		for y := -ey; y < scr.vsc.data.GetHeight(); y += bg.cell_height {
			for x := -ex; x < scr.vsc.data.GetWidth(); x += bg.cell_width {
				n := ((x + ex + dx) / bg.cell_width % bg.row) + ((y+ey+dy)/bg.cell_height%bg.column)*bg.row
				num := bg.data[n]
				if num < len(bg.pat) || num != 0 {
					pat := bg.pat[num]
					scr.vsc.data.PutPixelPat(x, y, pat)
				}
			}
		}
	}

}
func caldiff(a int, ca int, n int) (int, int) {
	w := a * ca
	dx := ((0-n)%w + w) % w
	ex := dx % ca
	dx = dx / ca * ca
	return dx, ex
}

func (scr *SpriteScreen) DrawLowest() {
	scr.vsc.data.PutPixelPat(0, 0, scr.back)
}

func (scr *SpriteScreen) CheckSpriteHit(id int) []int {
	sp := scr.sp[id-1]
	sp_x := sp.pos.X - sp.pat.center.X
	sp_y := sp.pos.Y - sp.pat.center.Y
	checks := make(map[int]bool)
	hits := make([]int, 0)
	for y := sp_y; y < sp.pat.pix.GetHeight(); y++ {
		for x := sp_x; x < sp.pat.pix.GetWidth(); x++ {
			id := scr.vsc.GetAttr(x, y)
			if id != 0 && id != sp.id && !checks[id] {
				checks[id] = true
				hits = append(hits, id)
			}
		}
	}
	return hits
}

func NewVscreen(w int, h int) *Vscreen {
	vscr := new(Vscreen)
	vscr.attr = NewAttributeScreen(w, h)
	return vscr
}

func (vsc *Vscreen) SetAttr(x int, y int, id int) {
	vsc.attr[x+y*vsc.data.GetWidth()] = id
}

func (vsc *Vscreen) GetAttr(x int, y int) int {
	return vsc.attr[x+y*vsc.data.GetWidth()]
}

type AttributeScreen []int

func NewAttributeScreen(w int, h int) (attr AttributeScreen) {
	attr = make([]int, w*h)
	return attr
}
