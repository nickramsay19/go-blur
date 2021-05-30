package main

import (
	"fmt"
	"os"
	"strconv"
	"image"
	"image/color"
	"github.com/lucasb-eyer/go-colorful"
)

type Pixel struct {
	X, Y int
}

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

func GetPixelsInRadius(pixels [][]color.RGBA64, x, y, dx, dy, r int) []color.RGBA64 {

	// create slice of adjacent pixels
	var adjacentPixels []color.RGBA64
	width, height := len(pixels), len(pixels[0])

	// check for recursive base case: repeats = 0
	if r == 0 {
		return adjacentPixels
	}

	// loop through all adjacent x + dx, y + dy values rdx,rdy within 1 pixel of x + dx,y + dy
	for rdx := -1; rdx < 1; rdx++ {
		for rdy := -1; rdy < 1; rdy++ {

			// create some 2d Pixel values to help track & describe the following operations
			callingPixel := Pixel {x, y}
			originPixel := Pixel {x + dx, y + dy}
			originPixelAdjacent := Pixel {x + dx + rdx, y + dy + rdy}

			// avoid adding this pixel (at x + dx, y + dy) 
			// avoid adding the calling pixel (at x, y)
			// avoid adding pixels adjacent to the calling pixel
			if (originPixelAdjacent.X != originPixel.X && originPixelAdjacent.Y != originPixel.Y) && (originPixelAdjacent.X != callingPixel.X && originPixelAdjacent.Y != callingPixel.Y) {

				// check that rdx and rdy make x,y are contained within the bounds of the slice
				if originPixelAdjacent.X > 0 && originPixelAdjacent.X < width - 1 && originPixelAdjacent.Y > 0 && originPixelAdjacent.Y < height - 1 {

					// append the adjacent pixel
					adjacentPixels = append(adjacentPixels, pixels[originPixelAdjacent.X][originPixelAdjacent.Y])

					// append further adjacent pixels recursively
					adjacentPixels = append(adjacentPixels, GetPixelsInRadius(pixels, originPixel.X, originPixel.Y, rdx, rdy, r - 1)...)
				}
			}
		}
	}

	// remove duplicates
	// duplicates occur when function called recursively adds the adjacent pixel that called it
	// however, duplicates also occur when two colors happen to be the same, in this case we don't want to remove those duplicates
	// so we wont for now
	// we could remove all adjacent pixels found in this running from the ones found in the recursively ran


	return adjacentPixels
}

func BlurPixel(pixels [][]color.RGBA64, x, y int, blurFactor float64, blurRadius int) color.RGBA64 {

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
	var adjacentPixels []color.RGBA64 = GetPixelsInRadius(pixels, x, y, 0, 0, blurRadius)
	
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
	//fmt.Printf("Blurred (%d, %d)\n", x, y)
	return blurredPixel
}

func main() {

	// Get command line arguments
	args := os.Args[1:]

	// declare blur params
	var inputImageFileName string
	var outputImageFileName string
	var blurFactor float64
	var blurRadius int

	// check that sufficient arguments have been provided
	// initialise parameters from args
	if len(args) < 2 {

		// insufficient args provided, show an error and usage guide.
		fmt.Println("Error: Insufficient parameters provided.\nUsage: blur <input> <output> # use default blur parameters\nUsage: blur <input> <output> <blur factor> <blur radius> # use specified blur parameters")
		return
	} else if len(args) < 4 {
		inputImageFileName = args[0]
		outputImageFileName = args[1]

		// give default values for blurFactor and blurRadius
		blurFactor = 0.1
		blurRadius = 2
	} else {
		inputImageFileName = args[0]
		outputImageFileName = args[1]

		// convert 3rd param to float
		// assign blurFactor to float if successful
		if blurFactorFloat, err := strconv.ParseFloat(args[2], 64); err == nil {

			// check for a valid float value between 0 and 1
			if blurFactorFloat >= 0.0 && blurFactorFloat <= 1 {
				blurFactor = blurFactorFloat
			} else {
				fmt.Println("Error: Improper value for blur factor provided.\nPlease provide a floating point number between 0 and 1.\nExample: 0.4\n")
			}
			
		} else {
			fmt.Println("Error: Improper value for blur factor provided.\nPlease provide a floating point number.\nExample: 0.4\n")
			return
		}

		// convert 4th param to int
		// assing to blurRadius if successful
		if blurRadiusFloat, err := strconv.Atoi(args[3]); err == nil {

			// check for a valid integer value between 0 and 10
			if blurRadiusFloat >= 0 && blurRadiusFloat <= 10 {
				blurRadius = blurRadiusFloat
			} else {
				fmt.Println("Error: Improper value for blur radius provided.\nPlease provide an integer between 0 and 10.\nExample: 2\n")
			}
			
		} else {
			fmt.Println("Error: Improper value for blur radius provided.\nPlease provide an integer.\nExample: 2\n")
			return
		}
	}

	fmt.Println("Blurring image.")

	// create and open the input image
	pixelImage := NewPixelImage(inputImageFileName, 810, 577)
	pixelImage = pixelImage.ReadFromFile()

	// create the result image from the blurred input image
	blurredPixelImage := pixelImage.Blurred(outputImageFileName, blurFactor, blurRadius)

	// write the blurred image to a jpeg
	blurredPixelImage.WriteToFile()

	// print a success message
	fmt.Println("Job finished successfully.")
}