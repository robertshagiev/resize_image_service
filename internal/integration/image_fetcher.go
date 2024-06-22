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

type imageFetcher struct {
	client *http.Client
}

func NewImageFetcher(timeout time.Duration) *imageFetcher {
	client := &http.Client{
		Timeout: timeout,
	}
	return &imageFetcher{client: client}
}

func (f *imageFetcher) FetchImage(url string) (*model.ImageData, error) {
	ctx := context.Background() // Используем простой контекст

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

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

	return &model.ImageData{
		Data:        data,
		ContentType: resp.Header.Get("Content-Type"),
		URL:         url,
		Width:       img.Width,
		Height:      img.Height,
	}, nil
}
