package main
import (
	"fmt"
	"os"
	_"image"
	"image/png"
	"strings"
)

func main () {

	args := rm(os.Args, 0)

	if !validateArgs(args) {
		return
	}

	for _, v := range args {
		file, err := os.Open(v)
		if err != nil {
			fmt.Printf("There was an error opening a file.\n")
			return
		}
		defer file.Close()

		img, err := png.Decode(file) 
		if err != nil {
			fmt.Printf("There was an error decoding a file. %v\n", err)
			return
		}

		bounds := img.Bounds()

		newfile, err := os.OpenFile(v[0:len(v)-4]+".jacob", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("There was an error creating the new file.\n")
			return
		}
		
		boundsStr := fmt.Sprintf("%d~%d|", bounds.Max.X, bounds.Max.Y) 
		_, err = newfile.Write([]byte(boundsStr))
				if err != nil {
					fmt.Printf("Error writing to file.")
					return
				}

		fmt.Printf("Converting %v...\n", v)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x,y).RGBA()

				r /= 256
				g /= 256
				b /= 256

				str := fmt.Sprintf("%d", r) + "/" + fmt.Sprintf("%d", g) + "/" + fmt.Sprintf("%d", b) + "|"

				_, err := newfile.Write([]byte(str))
				if err != nil {
					fmt.Printf("Error writing to file.")
					return
				}
			}
		}

		defer newfile.Close()
	}

}

func validateArgs(args []string) bool {

	// Check that an argument is provided
	if len(args) < 1 {
		fmt.Printf("No argument provided.\n")
		return false;
	}

	// Check last 4 chars of each argument to see if they are valid PNGs

	var valid bool
	for _, v := range args {
		length := len(v)

		if length < 5 {
			valid = false
			break;
		}

		var extension string
		extension = v[strings.LastIndex(v, "."):]

		if extension != ".png" {
			valid = false
			break;
		}

		valid = true
	}

	if !valid {
		fmt.Printf("One or more of the provided files are not .png.\n")
		return false
	}

	return true
	
}

func rm (args []string, index int) []string {
	return append(args[:index], args[index+1:]...)
}