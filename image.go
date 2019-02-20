package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type record struct {
	Description string
	Price       string
	imagePath   string
}

func readCsv(fileName string) []record {
	records := make([]record, 0)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error: %v", err)
		}
		r := record{
			Description: line[0],
			Price:       line[1],
			imagePath:   line[2],
		}
		records = append(records, r)
	}
	return records
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Error: %v", r)
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the name of the csv file: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	csvFileName := strings.TrimSpace(input)
	ImageDirectory := FindDir("Images")
	if !strings.Contains(ImageDirectory, "Images") {
		log.Fatalf("Images folder not found! Please check it.")
	}
	records := readCsv(csvFileName)
	for _, record := range records {
		fmt.Println(record)
	}
	// im, err := gg.LoadImage("./Images/2.jpg")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// b := im.Bounds()

	// S := float64(b.Max.X)
	// H := float64(b.Max.Y)

	// dc := gg.NewContext(int(S), int(H)+700)
	// dc.SetRGB(1, 1, 1)
	// dc.Clear()
	// dc.SetRGB(0, 0, 0)
	// if err := dc.LoadFontFace("./Fonts/Times New Roman.ttf", 180); err != nil {
	// 	panic(err)
	// }
	// dc.DrawRoundedRectangle(0, 0, S, H, 0)
	// dc.DrawImage(im, 0, 0)
	// dc.DrawStringAnchored("Ref.2430  Dia. 3.43 ct, Ruby 2 ct", S/2, H+170, 0.5, 0.5)
	// dc.DrawStringAnchored("13500 CHF", S/2, H+400, 0.5, 0.5)

	// dc.Clip()
	// dc.SavePNG("1.jpg")
}

func FindDir(dir string) string {
	// setting directory to the current
	fileName := "."
	// Checking if it is in the current directory  "./"
	if _, err := os.Stat("./" + dir + "/"); err == nil {
		// setting filepath to the the current directory
		fileName, _ = filepath.Abs("./" + dir + "/")
		// Checking if it is in the previous directory  "../"
	} else if _, err := os.Stat("../" + dir + "/"); err == nil {
		// setting filepath to the the previous directory
		fileName, _ = filepath.Abs("../" + dir + "/")
	}
	// returning absolute file path
	return fileName + "/"
}
