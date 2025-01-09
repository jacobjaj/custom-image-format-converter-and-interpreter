package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"image"
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
)

func main() {

	args := os.Args

	if !validateArgs(args) {
		return
	}

	a := app.New()
	window := a.NewWindow("Jacob Image Viewer")

	content, err := processFile(args[1])
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}

	sizeParams := strings.Split(content[0], "~")

	w, _ := strconv.Atoi(sizeParams[0])
	h, _ := strconv.Atoi(sizeParams[1])

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	var count int64 = 1
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			colors := splitColor(content[count])
			img.Set(x, y, color.RGBA{uint8(colors[0]), uint8(colors[1]), uint8(colors[2]), 255})
			count++
		}
	}

	raster := canvas.NewImageFromImage(img)

	window.SetContent(raster)
	window.Resize(fyne.NewSize(float32(w), float32(h)))
	window.ShowAndRun()
}

func splitColor(str string) [3]int {
	colors := strings.Split(str, "/")

	var out [3]int
	out[0],_ = strconv.Atoi(colors[0])
	out[1], _ = strconv.Atoi(colors[1])
	out[2], _ = strconv.Atoi(colors[2])

	return out
	
}

func processFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
			return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	var allFields []string // Single slice to hold all fields

	scanner := bufio.NewScanner(file)

	const maxFileSize = 20 * 1024 * 1024 // 20MB - adjust as needed
	buf := make([]byte, maxFileSize)
	scanner.Buffer(buf, maxFileSize)

	for scanner.Scan() {
			line := scanner.Text()
			fields := splitLine(line)

			// Append all fields from the current line to the main slice
			allFields = append(allFields, fields...) // Use ... to append elements of a slice
	}

	if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("scanning file: %w", err)
	}

	return allFields, nil
}

func splitLine(line string) []string {
	return strings.Split(line, "|")
}

func validateArgs(args []string) bool {

	// Check that an argument is provided
	if len(args) < 2 {
		fmt.Printf("No argument provided.\n")
		return false;
	}

	if len(args) > 2 {
		fmt.Printf("Too many arguments.\n")
		return false;
	}

	arg := args[1]

	// Check last 4 chars of each argument to see if they are valid PNGs

	length := len(arg)

	if length < 7 {
		fmt.Printf("The provided file is not a .jacob.\n")
		return false
	}

	var extension string
	extension = arg[strings.LastIndex(arg, "."):]

	if extension != ".jacob" {
		fmt.Printf("The provided file is not a .jacob.\n")
		return false
	}

	return true
	
}

