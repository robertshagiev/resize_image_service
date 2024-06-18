package integration

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"resize_image_service/internal/model"
)

type Integration struct{}

func NewIntegration() *Integration {
	return &Integration{}
}

func (integration *Integration) GetImageApiReq(url string) (*model.ImageData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	img, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return &model.ImageData{
		Data:        data,
		ContentType: resp.Header.Get("Content-Type"),
		URL:         url,
		Width:       img.Width,
		Height:      img.Height,
	}, nil
}
