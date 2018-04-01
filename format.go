package resize

import (
	"image"
	"os"
	"net/http"
	"io"
	"io/ioutil"
	"bytes"
)

type fileImage struct {
	filePath string
	mimeType string
	img image.Image
	bytesReader io.Reader
}

func LoadFile(filePath string) (f fileImage) {
	file := fileImage{filePath, "", nil, nil}
	file.UpdateMimeType()

	return file
}

func (f *fileImage) UpdateByteReader() (error) {
	file, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		return err
	}

	f.bytesReader = bytes.NewReader(file)

	return nil
}

func (f *fileImage) UpdateMimeType() (error) {

	osFile, err := os.Open(f.filePath)
	if err != nil {
		return err
	}
	defer osFile.Close()

	buffer := make([]byte, 512)

	_, er := osFile.Read(buffer)
	if er != nil {
		return er
	}

	contentType := http.DetectContentType(buffer)

	f.mimeType = contentType
	return nil
}