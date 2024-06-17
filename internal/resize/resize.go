package resize

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/disintegration/imaging"
)

type Resizer struct{}

func NewResizer() *Resizer {
	return &Resizer{}
}

func (r *Resizer) ResizeImage(url string, width, height int) ([]byte, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		return nil, "", err
	}

	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, resizedImg, nil)
	case "png":
		err = png.Encode(&buf, resizedImg)
	default:
		err = jpeg.Encode(&buf, resizedImg, nil)
	}
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), format, nil
}
