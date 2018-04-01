package resize

import (
	"image"
	"math"
	"image/color"
)

func NearestNeighbor(img image.Image, width uint, height uint) (image.Image) {

	pixelsX := float32(img.Bounds().Max.X) / float32(width)
	pixelsY := float32(img.Bounds().Max.Y) / float32(height)
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	for i := 0 ; i < int(width); i++ {
		for k := 0 ; k < int(height); k++ {
			pixelWidth := float32(i) * pixelsX
			pixelHeight := float32(k) * pixelsY

			if pixelWidth > float32(img.Bounds().Max.X) {
				pixelWidth = float32(img.Bounds().Max.X)
			}
			if pixelHeight > float32(img.Bounds().Max.Y) {
				pixelHeight = float32(img.Bounds().Max.Y)
			}

			red, green, blue, alpha := img.At(int(pixelWidth), int(pixelHeight)).RGBA()

			newImage.Set(i, k, color.RGBA{uint8(red / 257), uint8(green / 257), uint8(blue / 257), uint8(alpha / 257)})
		}
	}

	return newImage
}

func Supersample(img image.Image, width uint, height uint) (image.Image){

	pixelsX := float32(img.Bounds().Max.X) / float32(width)
	pixelsY := float32(img.Bounds().Max.Y) / float32(height)
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	for i := 0 ; i < int(width); i++ {
		for k := 0 ; k < int(height); k++ {
			pixelWidth := float32(i) * pixelsX
			pixelHeight := float32(k) * pixelsY

			if pixelWidth > float32(img.Bounds().Max.X) {
				pixelWidth = float32(img.Bounds().Max.X)
			}
			if pixelHeight > float32(img.Bounds().Max.Y) {
				pixelHeight = float32(img.Bounds().Max.Y)
			}

			var r []uint32
			var g []uint32
			var b []uint32
			var a []uint32

			radiusX := float32(math.Ceil(float64(pixelsX))) / 2
			radiusY := float32(math.Ceil(float64(pixelsY))) / 2

			for ii := pixelWidth - radiusX ; ii < pixelWidth + radiusX + 1 ; ii++ {
				for kk := pixelHeight - radiusY ; kk < pixelHeight + radiusY + 1 ; kk++ {
					pw := ii
					ph := kk
					if pw < 0 {
						pw = 0
					}
					if pw > float32(img.Bounds().Max.X) {
						pw = float32(img.Bounds().Max.X)
					}
					if ph < 0 {
						ph = 0
					}
					if ph > float32(img.Bounds().Max.Y) {
						ph = float32(img.Bounds().Max.Y)
					}

					ri, gi, bi, ai := img.At(int(pw), int(ph)).RGBA()
					r = append(r, ri)
					g = append(g, gi)
					b = append(b, bi)
					a = append(a, ai)
				}
			}

			red := sum(r) / uint32(len(r))
			green := sum(g) / uint32(len(g))
			blue := sum(b) / uint32(len(b))
			alpha := sum(a) / uint32(len(a))

			newImage.Set(i, k, color.RGBA{uint8(red / 257), uint8(green / 257), uint8(blue / 257), uint8(alpha / 257)})
		}
	}

	return newImage
}

func sum(input []uint32) uint32 {
	var sum uint32 = 0

	for _, i := range input {
		sum += i
	}

	return sum
}
