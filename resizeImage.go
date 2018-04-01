package resize

import (
	"image"
	"io/ioutil"
	"image/jpeg"
	"bytes"
	"os"
	"net/http"
	"image/png"
	"image/gif"
	"errors"
	"fmt"
)

// Типы изображений
const JPEG = "image/jpeg"
const PNG = "image/png"
const GIF = "image/gif"

// Алгоритмы преобразования
const SUPERSAMPLE = 1
const NEAREST_NEIGHBOR = 2

// Тип вписывания изображения
const COVER = 1
const DISTORTON = 2

type Resizer struct {
	algorithm int
	inscribe int
	mimeType string
	img image.Image
	newImg image.Image
}

func GetResizer() (r Resizer){
	return Resizer{
		1,
		1,
		"",
		nil,
		nil}
}

func (r *Resizer) SetAlgorithm(algorithm int) {
	r.algorithm = algorithm
}

func (r *Resizer) SetInscribe(inscribe int) {
	r.inscribe = inscribe
}

func (r *Resizer) Load(filePath string) (err error){

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	osFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer osFile.Close()

	mimeType, err := r.GetFileContentType(osFile)
	if err != nil{
		return err
	}
	r.mimeType = mimeType

	switch mimeType {
	case JPEG:
		r.img, err = jpeg.Decode(bytes.NewReader(file))
		return err
	case PNG:
		r.img, err = png.Decode(bytes.NewReader(file))
		return err
	case GIF:
		r.img, err = gif.Decode(bytes.NewReader(file))
		return err
	default:
		return errors.New(fmt.Sprintf("Mime type %s not supported", mimeType))
	}
}

func (r *Resizer) Resize(width uint, height uint) {

	switch r.inscribe {
	case COVER:
		width, height = r.Cover(width, height)
	}

	switch r.algorithm {
	case NEAREST_NEIGHBOR:
		r.newImg = r.NearestNeighbor(width, height)
	default:
		r.newImg = r.Supersample(width, height)
	}
}

func (r *Resizer) Cover(width uint, height uint) (newWidth uint, newHeight uint){
	koefOld := float32(r.img.Bounds().Max.X) / float32(r.img.Bounds().Max.Y)
	koefNew := float32(width) / float32(height)
	switch true {
	case width == 0:
		return uint(float32(height) * koefOld), height
	case height == 0:
		return width, uint(float32(width) / koefOld)
	}

	switch true {
	case koefNew > koefOld:
		return uint(float32(height) * koefOld), height
	case koefNew < koefOld:
		return width, uint(float32(width) / koefOld)
	default:
		return width, height
	}
}

func (r *Resizer) GetFileContentType(out *os.File) (string, error) {

	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}