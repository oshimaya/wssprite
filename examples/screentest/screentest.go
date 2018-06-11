package main

import (
	"github.com/oshimaya/wssprite"
	"image"
	"image/color"
	"time"
)

func main() {
	scr, err := wssprite.NewSpriteScreen("/dev/ttyE1", 256, 144, 1, 1)
	if err != nil {
		return
	}
	defer scr.Close()

	scr.CreateBg(20, 15, 16, 16)

	sp := scr.GetSprite(1)

	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.NRGBA{0, 255, 255, 255})
			} else {
				img.Set(x, y, color.NRGBA{0, 255, 255, 0})
			}
		}
	}
	pat := scr.CreatePatternFromImage(img)
	sp.SetPattern(pat)
	sp.SetPosition(160, 120)
	sp.Enable()
	img = image.NewRGBA(image.Rect(0, 0, 16, 16))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.NRGBA{0, 255, 0, 255})
		}
	}
	bgpat := scr.CreatePatternFromImage(img)

	bg := scr.GetBg(1)
	bg.AddPixelPattern(bgpat.GetPix())
	for y := 0; y < 30; y++ {
		for x := 0; x < 40; x++ {
			if (x+y)%2 == 0 {
				bg.Put(x, y, 1)
			}
		}
	}
	//	bg.SetScrollable(false)

	for i := 0; i < 200; i++ {
		bg.SetOffset(0-i, 0)
		sp.SetPosition(160-i, 120-i/2)
		scr.DrawLowest()
		scr.DrawBg(1)
		scr.DrawAllSprite()
		scr.DrawVscreen()
		time.Sleep(time.Millisecond * 16)
	}
	time.Sleep(time.Second * 5)
}
