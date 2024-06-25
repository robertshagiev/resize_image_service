package integration

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"resize_image_service/internal/model"
	"time"
)

type ImageFetcher struct {
	client  clientHttp
	timeout time.Duration
}

//go:generate  mockery --with-expecter --name clientHttp
type clientHttp interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewImageFetcher(timeout time.Duration) *ImageFetcher {
	client := &http.Client{
		Timeout: timeout,
	}
	return &ImageFetcher{
		client:  client,
		timeout: timeout,
	}
}

func (f *ImageFetcher) FetchImage(ctx context.Context, url string) (*model.ImageData, error) {
	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	fmt.Println(req)
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(data))
	}

	img, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	fmt.Println(img.Height, img.Width)
	fmt.Println(data)
	fmt.Println(resp.Header.Get("Content-Type"))
	return &model.ImageData{
		Data:        data,
		ContentType: resp.Header.Get("Content-Type"),
		URL:         url,
		Width:       img.Width,
		Height:      img.Height,
	}, nil
}
