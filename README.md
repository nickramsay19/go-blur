# go-blur
> Created by Nicholas Ramsay

go-blur is a command line utility for blurring JPEG images. 

## Usage
```
# build
go build main.go PixelImage.go

# run on your image
# ./main <input> <output> <blur factor> <blur radius>
./main myInputImage.jpg myOutputImage.jpg 0.1 2
```

The blur radius specifies the radius in which each pixel will be blended with. The lower the value the faster the program runs. I recommend a value of 2 or 3.

The blur factor specifies the amount that a pixel will be blended with surrounding pixels, the value describes the amount of the original color will remain. I recommend a value of 0.1.

## Roadmap
- [x] Refactor and clean the code
- [x] Radii blurring
- [x] Allow command line inputs
- [ ] Use concurrency
- [ ] Implement my own color blending algorithm
- [x] Refactor w/ PixelImage struct
- [x] GetPixelInRadius as a for loop -1 - +1
- [ ] Check for valid filenames as given from args