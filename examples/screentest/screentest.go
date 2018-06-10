package main

import (
	"github.com/oshimaya/wssprite"
	"image"
	"image/color"
	"time"
)

func main() {
	scr, err := wssprite.NewSpriteScreen("/dev/ttyE1", 320, 240, 1, 1)
	if err != nil {
		return
	}
	defer scr.Close()

	scr.CreateBg(40, 30, 8, 8)

	sp := scr.GetSprite(1)

	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, color.NRGBA{0, 255, 255, 255})
		}
	}
	pat := scr.CreatePatternFromImage(img)
	sp.SetPattern(pat)
	sp.SetPosition(160, 120)
	sp.Enable()
	for i := 0; i < 100; i++ {
		sp.SetPosition(160-i, 120)
		scr.DrawAllSprite()
		scr.DrawVscreen()
		time.Sleep(time.Millisecond * 16)
	}
	time.Sleep(time.Second * 5)
}
