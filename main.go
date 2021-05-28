package main

import (
	"fmt"
	"os"
	"image"
	//"image/draw"
	"image/color"
	"image/jpeg"
	"github.com/lucasb-eyer/go-colorful"
)

func getRGBA64At(img image.Image, x, y int) color.RGBA64 {

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

func main() {
	fmt.Println("heyyyyy")

	// open the image
	fimg, _ := os.Open("img.jpg")
	defer fimg.Close()
	img, _, _ := image.Decode(fimg)

	// copy image into pixel array
	var pixels [225][225]color.RGBA64

	// loop through all x and y positions of the image
	for x := 0; x < 225; x++ {
		for y := 0; y < 225; y++ {
			pixels[x][y] = getRGBA64At(img, x, y)
		}
	}

	// make a copy of pixels for blurring
	var blurredPixels [225][225]color.RGBA64
	var blurFactor float64 = 0.1
	for x := 0; x < 225; x++ {
		for y := 0; y < 225; y++ {
			
			// we now begin the blurring process of the pixel at x,y
			// get original pixel color
			originalPixel := getRGBA64At(img, x, y)

			// create slice of adjacent pixels
			var adjacentPixels []color.RGBA64

			// find all adjacent pixels and append to slice
			if x > 0 {
				// there are pixels left of x,y
				adjacentPixels = append(adjacentPixels, getRGBA64At(img, x - 1, y))
				if y > 0 {
					// there are pixels above x,y
					adjacentPixels = append(adjacentPixels, getRGBA64At(img, x - 1, y - 1))
				}
				if y < 255 {
					// there are pixels below x,y
					adjacentPixels = append(adjacentPixels, getRGBA64At(img, x - 1, y + 1))
				}
			}
			if x < 255 {
				// there are pixels right of x,y
				adjacentPixels = append(adjacentPixels, getRGBA64At(img, x + 1, y))
				if y > 0 {
					// there are pixels above x,y
					adjacentPixels = append(adjacentPixels, getRGBA64At(img, x + 1, y - 1))
				}
				if y < 255 {
					// there are pixels below x,y
					adjacentPixels = append(adjacentPixels, getRGBA64At(img, x + 1, y + 1))
				}
			}
			if y > 0 {
				// there is a pixel above x,y
				adjacentPixels = append(adjacentPixels, getRGBA64At(img, x, y - 1))
			}
			if y < 255 {
				// there is a pixel below x,y
				adjacentPixels = append(adjacentPixels, getRGBA64At(img, x, y + 1))
			}


			// create a colorful Color from the normalised original pixel vals
			colorfulOriginalPixel := colorful.Color{
				float64(originalPixel.R) / 65535, 
				float64(originalPixel.G) / 65535, 
				float64(originalPixel.B) / 65535,
			}
		
			// now set blurredPixels at x,y to the blurred pixel
			//blurredPixels[x][y] = blurredPixel
			//continue

			// get the average of the adjacent colors
			// do this by finding the total of each color and then dividing by the amount of colors
			// we use uint64 for this part of the operation to fit larger numbers
			//var totalAdjacentPixelsRed, totalAdjacentPixelsGreen, totalAdjacentPixelsBlue uint64 = 0, 0, 0
			
			for _, p := range adjacentPixels {
				colorfulAdjacentPixel := colorful.Color{
					float64(p.R) / 65535, 
					float64(p.G) / 65535, 
					float64(p.B) / 65535,
				}
				colorfulOriginalPixel = colorfulOriginalPixel.BlendRgb(colorfulAdjacentPixel, blurFactor)
			}

			// declare and initialise the average of the adjacent pixels
			//averageAdjacentPixel := color.RGBA64{
			//	uint16(totalAdjacentPixelsRed) / uint16(len(adjacentPixels)),
			//	uint16(totalAdjacentPixelsGreen) / uint16(len(adjacentPixels)),
			//	uint16(totalAdjacentPixelsBlue) / uint16(len(adjacentPixels)),
			//	65535,
			//}

			// now find the average between the original and the averageAdjacentPixel given a blurring factor
			var blurredPixel color.RGBA64 = color.RGBA64{
				uint16(colorfulOriginalPixel.R * 65535),
				uint16(colorfulOriginalPixel.G * 65535),
				uint16(colorfulOriginalPixel.B * 65535),
				65535,
			}
		
			// now set blurredPixels at x,y to the blurred pixel
			blurredPixels[x][y] = blurredPixel
		}
	}

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