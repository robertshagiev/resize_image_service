package integration

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"resize_image_service/internal/model"

	"github.com/go-resty/resty/v2"
)

type imageFetcher struct {
	client *resty.Client
}

func NewImageFetcher(client *resty.Client) *imageFetcher {
	return &imageFetcher{client: client}
}

func (f *imageFetcher) FetchImage(url string) (*model.ImageData, error) {
	ctx := context.Background()
	req, err := NewRequestContext(ctx, f.client, url)
	if err != nil {
		return nil, err
	}

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	data := resp.Body()
	img, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return &model.ImageData{
		Data:        data,
		ContentType: resp.Header().Get("Content-Type"),
		URL:         url,
		Width:       img.Width,
		Height:      img.Height,
	}, nil
}
