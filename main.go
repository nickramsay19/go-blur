package main

import (
	"fmt"
	"os"
	"sync"
	"image"
	//"image/draw"
	"image/color"
	"image/jpeg"
	"github.com/lucasb-eyer/go-colorful"
)

func GetRGBA64At(img image.Image, x, y int) color.RGBA64 {

	// declare the output rgba64
	var rgba64 color.RGBA64

	// declare and initialise temporary rgba values to pixel at x,y
	// we need temporary values to store the rgba vals as uint32
	var r, g, b, a uint32 = img.At(x, y).RGBA()

	// assign pixel values from temporary rgba values and convert to uint16 needed by color.RGBA64
	rgba64.R = uint16(r)
	rgba64.G = uint16(g)
	rgba64.B = uint16(b)
	rgba64.A = uint16(a)

	// finally return the produced color
	return rgba64
}

func GetPixelsInRadius(pixels [225][225]color.RGBA64, x, y, r int) []color.RGBA64 {

	// create slice of adjacent pixels
	var adjacentPixels []color.RGBA64

	if r == 0 {
		return adjacentPixels
	}

	// find all adjacent pixels and append to slice
	if x > 0 {
		// there are pixels left of x,y
		adjacentPixels = append(adjacentPixels, pixels[x - 1][y])
		adjacentPixels = append(adjacentPixels, GetPixelsInRadius(pixels, x - 1, y, r - 1)...)
	} 
	if x > 0 && y > 0 {
		// there are pixels above x,y
		adjacentPixels = append(adjacentPixels, pixels[x - 1][y - 1])
		adjacentPixels = append(adjacentPixels, GetPixelsInRadius(pixels, x - 1, y - 1, r - 1)...)
	} 
	if x > 0 && y < 224 {
		// there are pixels below x,y
		adjacentPixels = append(adjacentPixels, pixels[x - 1][y + 1])
		adjacentPixels = append(adjacentPixels, GetPixelsInRadius(pixels, x - 1, y + 1, r - 1)...)
	} 
	if x < 224 {
		// there are pixels right of x,y
		adjacentPixels = append(adjacentPixels, pixels[x + 1][y])
		adjacentPixels = append(adjacentPixels, GetPixelsInRadius(pixels, x + 1, y, r - 1)...)
	} 
	if x < 224 && y > 0 {
		// there are pixels above x,y
		adjacentPixels = append(adjacentPixels, pixels[x + 1][y - 1])
		adjacentPixels = append(adjacentPixels, GetPixelsInRadius(pixels, x + 1, y - 1, r - 1)...)
	} 
	if x < 224 && y < 224 {
		// there are pixels below x,y
		adjacentPixels = append(adjacentPixels, pixels[x + 1][y + 1])
		adjacentPixels = append(adjacentPixels, GetPixelsInRadius(pixels, x + 1, y + 1, r - 1)...)
	} 
	if y > 0 {
		// there is a pixel above x,y
		adjacentPixels = append(adjacentPixels, pixels[x][y - 1])
		adjacentPixels = append(adjacentPixels, GetPixelsInRadius(pixels, x, y - 1, r - 1)...)
	} 
	if y < 224 {
		// there is a pixel below x,y
		adjacentPixels = append(adjacentPixels, pixels[x][y + 1])
		adjacentPixels = append(adjacentPixels, GetPixelsInRadius(pixels, x, y + 1, r - 1)...)
	}

	// remove duplicates
	// duplicates occur when function called recursively adds the adjacent pixel that called it
	// however, duplicates also occur when two colors happen to be the same, in this case we don't want to remove those duplicates
	// so we wont for now
	// we could remove all adjacent pixels found in this running from the ones found in the recursively ran

	return adjacentPixels
}

func BlurPixel(pixels [225][225]color.RGBA64, x, y int, blurFactor float64, blurRadius int) color.RGBA64 {

	// we now begin the blurring process of the pixel at x,y
	// get original pixel color
	originalPixel := pixels[x][y]

	// create a pixel object as a colorful Color 
	// we need to normalise the RGB values for colorful color object
	colorfulOriginalPixel := colorful.Color {
		float64(originalPixel.R) / 65535, 
		float64(originalPixel.G) / 65535, 
		float64(originalPixel.B) / 65535,
	}

	// create slice of adjacent pixels
	var adjacentPixels []color.RGBA64 = GetPixelsInRadius(pixels, x, y, blurRadius)
	
	// loop through each discovered adjacent pixel
	// and mix the original pixel with the adjacent pixel
	for _, p := range adjacentPixels {

		// convert the adjacent pixel to a colorful color object
		// we must normalise the RGB vals for colorful color object
		colorfulAdjacentPixel := colorful.Color {
			float64(p.R) / 65535, 
			float64(p.G) / 65535, 
			float64(p.B) / 65535,
		}

		// we now blur the original pixel by mixing it in a colorful color lab with the adjacent pixel
		colorfulOriginalPixel = colorfulOriginalPixel.BlendRgb(colorfulAdjacentPixel, blurFactor)
	}

	// now find the average between the original and the averageAdjacentPixel given a blurring factor
	// we denormalise the RGB values as uint16
	var blurredPixel color.RGBA64 = color.RGBA64 {
		uint16(colorfulOriginalPixel.R * 65535),
		uint16(colorfulOriginalPixel.G * 65535),
		uint16(colorfulOriginalPixel.B * 65535),
		65535,
	}

	// finally return the blurred pixel
	fmt.Printf("Blurred (%d, %d)\n", x, y)
	return blurredPixel
}

func main() {
	fmt.Println("Blurring image.")

	var blurFactor float64 = 0.9
	var blurRadius int = 3

	// open the image
	fimg, _ := os.Open("img.jpg")
	defer fimg.Close()
	img, _, _ := image.Decode(fimg)

	// copy image into pixel array
	var pixels [225][225]color.RGBA64

	// loop through all x and y positions of the image
	for x := 0; x < 225; x++ {
		for y := 0; y < 225; y++ {

			// find the color at x,y and assign to pixels at x,y
			pixels[x][y] = GetRGBA64At(img, x, y)
		}
	}

	// make a copy of pixels for blurring
	var blurredPixels [225][225]color.RGBA64

	// loop through and blur each x and y pixel of the image
	var wg sync.WaitGroup
	for x := 0; x < 224; x++ {
		for y := 0; y < 224; y++ {
			
			wg.Add(1)

			// now set blurredPixels at x,y to the blurred pixel
			blurredPixels[x][y] = BlurPixel(pixels, x, y, blurFactor, blurRadius)
			wg.Done()
			//go func() {
				
			//}()
		}
	}

	wg.Wait()

	// create new blank image
	resImg := image.NewRGBA(image.Rect(0, 0, 225, 225))

	// set each pixel from the pixels array
	for x := 0; x < 225; x++ {
		for y := 0; y < 225; y++ {
			resImg.Set(x, y, blurredPixels[x][y])
		}
	}

	// create result image file
	resFImg, _ := os.Create("result.jpg")
	defer resFImg.Close()
	jpeg.Encode(resFImg, resImg, &jpeg.Options{jpeg.DefaultQuality})
}