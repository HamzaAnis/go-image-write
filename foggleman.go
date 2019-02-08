package main

import (
	"log"

	"github.com/fogleman/gg"
)

func main() {
	im, err := gg.LoadImage("./2.png")
	if err != nil {
		log.Fatal(err)
	}
	b := im.Bounds()

	S := float64(b.Max.X)
	H := float64(b.Max.Y)

	dc := gg.NewContext(int(S), int(H)+200)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("./Roboto-Regular.ttf", 30); err != nil {
		panic(err)
	}
	dc.DrawRoundedRectangle(0, 0, S, H, 0)
	dc.DrawImage(im, 0, 0)
	dc.DrawStringAnchored("Wow this is great", S/2, H+50, 0.5, 0.5)
	dc.DrawStringAnchored("13500 CHF", S/2, H+100, 0.5, 0.5)

	dc.Clip()
	dc.SavePNG("out1.png")
}
