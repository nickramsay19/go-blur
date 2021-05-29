package main

import (
	"fmt"
	"os"
	//"sync"
	"image"
	//"image/draw"
	"image/color"
	"image/jpeg"
	//"github.com/lucasb-eyer/go-colorful"
)

type PixelImage struct {
	imageFileName string
	pixels [][]color.RGBA64
	width, height int
}

func NewPixelImage(imageFileName string, width, height int) (newPixelImage PixelImage) {

	// initialise empty color slice
	var pixels [][]color.RGBA64 = make([][]color.RGBA64, width)

	// add each column to pixels
	for x := 0; x < width; x++ {
		pixels[x] = make([]color.RGBA64, height)
	}

	// create a blank PixelImage
	newPixelImage = PixelImage{imageFileName, pixels, width, height}

	// return produced pixelImage
	return
}

func (pi PixelImage) ReadFromFile() (newPixelImage PixelImage) {

	// open the image
	imgFile, _ := os.Open(pi.imageFileName)
	defer imgFile.Close()
	img, _, _ := image.Decode(imgFile)

	// copy image into pixel array
	pi.pixels = make([][]color.RGBA64, img.Bounds().Dx())
	fmt.Printf("during reading: %d\n", len(pi.pixels))

	// loop through all x and y positions of the image
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {

			// Check if pixels[x] exists
			// If not, initialise pixels[x] to a slice of color.RGBA64
			// This shouldn't occur since pixels for all x <= width have been initialised
			if len(pi.pixels) < x {
				pi.pixels = append(pi.pixels, make([]color.RGBA64, img.Bounds().Dy()))
			}

			// Check if pixels[x] is intialised
			// If not, initialise pixels[x] to a slice of color.RGBA64
			if pi.pixels[x] == nil {
				pi.pixels[x] = make([]color.RGBA64, img.Bounds().Dy())
			}

			// find the color at x,y and assign to pixels at x,y
			pi.pixels[x][y] = GetRGBA64At(img, x, y)
		}
	}
	fmt.Printf("*during reading: %d\n", len(pi.pixels))

	newPixelImage = PixelImage {pi.imageFileName, pi.pixels, img.Bounds().Dx(),  img.Bounds().Dy()}
	fmt.Printf("**during reading: %d\n", len(newPixelImage.pixels))
	return
}

func (pi PixelImage) WriteToFile() {

	// create new blank image
	img := image.NewRGBA(image.Rect(0, 0, 225, 225))

	// set each pixel from the pixels array
	for x := 0; x < 225; x++ {
		for y := 0; y < 225; y++ {
			img.Set(x, y, pi.pixels[x][y])
		}
	}

	// create result image file
	imgFile, _ := os.Create(pi.imageFileName)
	defer imgFile.Close()
	jpeg.Encode(imgFile, img, &jpeg.Options{jpeg.DefaultQuality})
}

func (pi PixelImage) Blurred(newImageFileName string, blurFactor float64, blurRadius int) (blurredPixelImage PixelImage) {
	blurredPixelImage = NewPixelImage(newImageFileName, pi.width, pi.height)
	
	for x := 0; x < pi.width; x++ {
		for y := 0; y < pi.height; y++ {
			
			// now set blurredPixels at x,y to the blurred pixel
			
			//go func() {
			blurredPixelImage.pixels[x][y] = BlurPixel(pi.pixels, x, y, blurFactor, blurRadius)
			//}()
		}
	}

	return
}