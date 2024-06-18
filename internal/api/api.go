package api

import (
	"io"
	"net/http"
	"resize_image_service/internal/model"
)

type Api struct{}

func NewAPI() *Api {
	return &Api{}
}

func (api *Api) GetImageApi(url string) (*model.ImageData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &model.ImageData{
		Data:        data,
		ContentType: resp.Header.Get("Content-Type"),
	}, nil
}
