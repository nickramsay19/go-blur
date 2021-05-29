package main

import (
	//"fmt"
	"os"
	//"sync"
	"image"
	//"image/draw"
	"image/color"
	"image/jpeg"
	//"github.com/lucasb-eyer/go-colorful"
)

// close
// create pixel image from file name

type PixelImage struct {
	imageFileName string
	pixels [][]color.RGBA64
}

func NewPixelImage(imageFileName string) {
	//var newPixelImage PixelImage {imageFileName}


	// open the image
	fimg, _ := os.Open(imageFileName)
	defer fimg.Close()
	img, _, _ := image.Decode(fimg)

	// copy image into pixel array
	var pixels [][]color.RGBA64 = make([][]color.RGBA64, img.Width)

	// loop through all x and y positions of the image
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {

			// Check if pixels[x] exists
			// If not, initialise pixels[x] to a slice of color.RGBA64
			// This shouldn't occur since pixels for all x <= width have been initialised
			if len(pixels) < x {
				pixels = append(pixels, make([]color.RGBA64, img.Height))
			}

			// Check if pixels[x] is intialised
			// If not, initialise pixels[x] to a slice of color.RGBA64
			if pixels[x] == nil {
				pixels[x] = make([]color.RGBA64, img.Height)
			}

			// find the color at x,y and assign to pixels at x,y
			pixels[x][y] = GetRGBA64At(img, x, y)
		}
	}
}