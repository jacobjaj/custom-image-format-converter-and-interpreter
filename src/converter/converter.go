package main
import (
	"fmt"
	"os"
	_"image"
	"io/ioutil"
	"image/png"
	"bytes"
	"encoding/binary"
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
		
		xBits :=  intToBits(bounds.Max.X)
		yBits := intToBits(bounds.Max.Y)

		writeBits := append(xBits, yBits...)

		fmt.Printf("Converting %v...\n", v)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x,y).RGBA()

				rByte := uint8(r >> 8)
				gByte := uint8(g >> 8)
				bByte := uint8(b >> 8)

				rgb := []byte{rByte, gByte, bByte}

				writeBits = append(writeBits, rgb...)
			}
		}

		err = ioutil.WriteFile(v[0:len(v)-4]+".jacob", writeBits, 0644)
		if err != nil {
			fmt.Printf("There was an error creating the new file.\n")
			return
		}
	}

}

func intToBits(i int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int16(i))
	if err != nil {
			panic(err)
	}
	return buf.Bytes()
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