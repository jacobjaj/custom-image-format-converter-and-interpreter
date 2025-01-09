package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"image"
	"fmt"
	"io"
	"os"
	"strings"
	"encoding/binary"
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

	w := int(binary.BigEndian.Uint16(content[0:2]))
	h := int(binary.BigEndian.Uint16(content[2:4]))

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	var count int64 = 4
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{content[count], content[count+1], content[count+2], 255})
			count += 3
		}
	}

	raster := canvas.NewImageFromImage(img)

	window.SetContent(raster)
	window.Resize(fyne.NewSize(float32(w), float32(h)))
	window.ShowAndRun()
}

func processFile(filePath string) ([]byte, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("opening file: %w", err)
    }
    defer file.Close()

    content, err := io.ReadAll(file) // Read the entire file at once
    if err != nil {
        return nil, fmt.Errorf("reading file: %w", err)
    }

    return content, nil
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

