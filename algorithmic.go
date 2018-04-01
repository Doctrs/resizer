package resize

import (
	"image"
	"math"
	"image/color"
)

func (r *Resizer) NearestNeighbor(width uint, height uint) (image.Image) {

	pixelsX := float32(r.img.Bounds().Max.X) / float32(width)
	pixelsY := float32(r.img.Bounds().Max.Y) / float32(height)
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	for i := 0 ; i < int(width); i++ {
		for k := 0 ; k < int(height); k++ {
			pixelWidth := float32(i) * pixelsX
			pixelHeight := float32(k) * pixelsY

			if pixelWidth > float32(r.img.Bounds().Max.X) {
				pixelWidth = float32(r.img.Bounds().Max.X)
			}
			if pixelHeight > float32(r.img.Bounds().Max.Y) {
				pixelHeight = float32(r.img.Bounds().Max.Y)
			}

			red, green, blue, alpha := r.img.At(int(pixelWidth), int(pixelHeight)).RGBA()

			newImage.Set(i, k, color.RGBA{uint8(red / 257), uint8(green / 257), uint8(blue / 257), uint8(alpha / 257)})
		}
	}

	return newImage
}

func (resizer *Resizer) Supersample(width uint, height uint) (image.Image){

	pixelsX := float32(resizer.img.Bounds().Max.X) / float32(width)
	pixelsY := float32(resizer.img.Bounds().Max.Y) / float32(height)
	newImage := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	radiusX := float32(math.Ceil(float64(pixelsX))) / 2
	radiusY := float32(math.Ceil(float64(pixelsY))) / 2

	arraySize := math.Ceil(float64((radiusY * 2 + 1) * (radiusX * 2 + 1)))

	var r []uint32
	var g []uint32
	var b []uint32
	var a []uint32

	for j := 0 ; j < int(arraySize) ; j++ {
		r = append(r, 0)
		g = append(g, 0)
		b = append(b, 0)
		a = append(a, 0)
	}

	for i := 0 ; i < int(width); i++ {
		for k := 0 ; k < int(height); k++ {
			pixelWidth := float32(i) * pixelsX
			pixelHeight := float32(k) * pixelsY

			if pixelWidth > float32(resizer.img.Bounds().Max.X) {
				pixelWidth = float32(resizer.img.Bounds().Max.X)
			}
			if pixelHeight > float32(resizer.img.Bounds().Max.Y) {
				pixelHeight = float32(resizer.img.Bounds().Max.Y)
			}


			count := 0
			for ii := pixelWidth - radiusX ; ii < pixelWidth + radiusX + 1 ; ii++ {
				for kk := pixelHeight - radiusY ; kk < pixelHeight + radiusY + 1 ; kk++ {
					pw := ii
					ph := kk
					if pw < 0 {
						pw = 0
					}
					if pw > float32(resizer.img.Bounds().Max.X) {
						pw = float32(resizer.img.Bounds().Max.X)
					}
					if ph < 0 {
						ph = 0
					}
					if ph > float32(resizer.img.Bounds().Max.Y) {
						ph = float32(resizer.img.Bounds().Max.Y)
					}

					ri, gi, bi, ai := resizer.img.At(int(pw), int(ph)).RGBA()

					r[count] = ri
					g[count] = gi
					b[count] = bi
					a[count] = ai

					count++
				}
			}

			red := sum(r) / uint32(arraySize)
			green := sum(g) / uint32(arraySize)
			blue := sum(b) / uint32(arraySize)
			alpha := sum(a) / uint32(arraySize)

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
