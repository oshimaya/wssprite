package wssprite

import (
	"github.com/oshimaya/gowsdisplay"
	"image"
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

	scr.vsc.data = pix
	sp := make([]Sprite, spnum)
	for i := range sp {
		sp[i] = NewSprite(i + 1)
	}
	scr.sp = sp
	scr.bg = make([]BgScreen, bgnum)
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
	return &scr.bg[num]
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
