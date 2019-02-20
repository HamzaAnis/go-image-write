package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the name of the csv file: ")
	fileName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fileName)
	im, err := gg.LoadImage("./Images/2.jpg")
	if err != nil {
		log.Fatal(err)
	}
	b := im.Bounds()

	S := float64(b.Max.X)
	H := float64(b.Max.Y)

	dc := gg.NewContext(int(S), int(H)+700)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("./Fonts/Times New Roman.ttf", 180); err != nil {
		panic(err)
	}
	dc.DrawRoundedRectangle(0, 0, S, H, 0)
	dc.DrawImage(im, 0, 0)
	dc.DrawStringAnchored("Ref.2430  Dia. 3.43 ct, Ruby 2 ct", S/2, H+170, 0.5, 0.5)
	dc.DrawStringAnchored("13500 CHF", S/2, H+400, 0.5, 0.5)

	dc.Clip()
	dc.SavePNG("1.jpg")
}
