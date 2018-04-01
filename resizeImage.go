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
}

func GetResizer() (r Resizer){
	return Resizer{
		1,
		1,
		""}
}

func (r *Resizer) SetAlgorithm(algorithm int) {
	r.algorithm = algorithm
}

func (r *Resizer) SetInscribe(inscribe int) {
	r.inscribe = inscribe
}

func (r *Resizer) Load(filePath string) (img image.Image, err error){

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	osFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer osFile.Close()

	mimeType, err := r.GetFileContentType(osFile)
	if err != nil{
		return nil, err
	}
	r.mimeType = mimeType

	switch mimeType {
	case JPEG:
		return jpeg.Decode(bytes.NewReader(file))
	case PNG:
		return png.Decode(bytes.NewReader(file))
	case GIF:
		return gif.Decode(bytes.NewReader(file))
	default:
		return nil, errors.New(fmt.Sprintf("Mime type %s not supported", mimeType))
	}
}

func (r *Resizer) Resize(img image.Image, width uint, height uint) image.Image {

	switch r.inscribe {
	case COVER:
		width, height = r.Cover(img, width, height)
	}

	switch r.algorithm {
	case NEAREST_NEIGHBOR:
		return NearestNeighbor(img, width, height)
	default:
		return Supersample(img, width, height)
	}
}

func (r *Resizer) Cover(img image.Image, width uint, height uint) (newWidth uint, newHeight uint){
	koefOld := float32(img.Bounds().Max.X) / float32(img.Bounds().Max.Y)
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