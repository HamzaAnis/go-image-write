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
func (r *record) toString(imageDir string){
	fmt.Printf("Processing\n\tFile name:   %v\n\tDescription: %v\n\tPrice:       %v\n",imageDir+r.imageName,r.Description,r.Price)
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

func processImage(imageDir string,data record) {
	defer wg.Done()
	lock <- 1
	data.toString(imageDir)
	<-lock
	im, err := gg.LoadImage(imageDir + data.imageName)
	if err != nil {
		log.Println(err)
		return
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

	outputPath:=filepath.FromSlash(imageDir+"Output/"+data.imageName)
	dc.Clip()
	dc.SavePNG("Output/"+data.imageName)

	lock <- 1
	fmt.Printf("%v saved to %v\n", data.imageName, outputPath)
	<-lock
	<-sem // removes an int from sem, allowing another to proceed

}

var wg sync.WaitGroup
const MaxThreads=20
var lock = make(chan int, 1)
// 20 threads at one time
var sem = make(chan int, MaxThreads)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error: %v", r)
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

	// creating output folder
	newpath := filepath.Join(".", "Output")
	os.MkdirAll(newpath, os.ModePerm)
	
	log.Printf("\nNumber of parallel processing of images = %v\n",MaxThreads)
	for i:=1;i<len(records);i++{
		wg.Add(1)
		sem <- 1 // will block if there is MAX ints in sem
		go processImage(ImageDirectory,records[i])
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
	return filepath.FromSlash(fileName + "/")
}
