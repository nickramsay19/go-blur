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

## Roadmap
- [x] Refactor and clean the code
- [x] Radii blurring
- [x] Allow command line inputs
- [ ] Use concurrency
- [ ] Implement my own color blending algorithm
- [x] Refactor w/ PixelImage struct
- [ ] GetPixelInRadius as a for loop -1 - +1