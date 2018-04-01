package resize

import (
	"image"
	"os"
	"image/jpeg"
	"image/draw"
	"image/color"
	"image/png"
	"errors"
	"fmt"
)

func (r *Resizer) ConvertImage(img image.Image) (newImage image.Image){
	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.FloydSteinberg.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min)

	return newImg
}

func (r *Resizer) Save(filePath string) (err error){
	switch r.mimeType {
	case JPEG:
		return r.SaveJpeg(filePath)
	case PNG:
		return r.SavePng(filePath)
	default:
		return errors.New(fmt.Sprintf("Mime type %s not supported", r.mimeType))
	}
}


func (r *Resizer) SaveJpeg(filePath string) (err error){
	if r.mimeType != JPEG {
		r.newImg = r.ConvertImage(r.newImg)
	}

	outputFile, err := os.Create(filePath)
	if err != nil{
		return err
	}

	jpeg.Encode(outputFile, r.newImg, &jpeg.Options{95})

	return nil
}

func (r *Resizer) SavePng(filePath string) (err error){
	if r.mimeType != PNG {
		r.newImg = r.ConvertImage(r.newImg)
	}

	outputFile, err := os.Create(filePath)
	if err != nil{
		return err
	}

	return png.Encode(outputFile, r.newImg)
}