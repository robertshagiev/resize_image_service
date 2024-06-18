package resize

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	"resize_image_service/internal/model"

	"github.com/disintegration/imaging"
)

type Resizer struct{}

func NewResizer() *Resizer {
	return &Resizer{}
}

func (r *Resizer) ResizeImage(data []byte, contentType string, width, height int) (*model.ResizedImage, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &model.ResizedImage{
		Data:   buf.Bytes(),
		Format: format,
	}, nil
}
