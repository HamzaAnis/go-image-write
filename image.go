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
	"sync"

	"github.com/fogleman/gg"
)

type record struct {
	Description string
	Price       string
	imageName   string
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
			imageName:   line[2],
		}
		records = append(records, r)
	}
	return records
}

func processImage(data record) {
	defer wg.Done()
	lock <- 1
	log.Printf("Processing %v with description %v and price %v", data.imageName, data.Description, data.Price)
	<-lock
	im, err := gg.LoadImage("./Images/" + data.imageName)
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
	dc.DrawStringAnchored(data.Description, S/2, H+170, 0.5, 0.5)
	dc.DrawStringAnchored(data.Price, S/2, H+400, 0.5, 0.5)

	dc.Clip()
	dc.SavePNG("./Output/" + data.imageName)

	lock <- 1
	log.Printf("Image/%v saved to Output/%v", data.imageName, data.imageName)
	<-lock
}

var wg sync.WaitGroup

var lock = make(chan int, 1)

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
		fmt.Println(record.imageName)
	}

	// creating output folder
	newpath := filepath.Join(".", "Output")
	os.MkdirAll(newpath, os.ModePerm)

	for _, record := range records[1:] {
		wg.Add(1)
		go processImage(record)
	}
	wg.Wait()
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
