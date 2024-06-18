package usecase

import (
	"resize_image_service/internal/model"
)

type ImageServiceIn interface {
	ResizeImage(url string, width, height int) (*model.ResizedImage, error)
}

type usecase struct {
	integration Integration
	resizer     ResizerIn
	logger      Log
}

type Integration interface {
	GetImageApiReq(url string) (*model.ImageData, error)
}

type ResizerIn interface {
	ResizeImage(data []byte, contentType string, width, height int) (*model.ResizedImage, error)
}

type Log interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}

func NewImageService(api Integration, resizer ResizerIn, log Log) *usecase {
	return &usecase{
		integration: api,
		resizer:     resizer,
		logger:      log,
	}
}

func (s *usecase) ResizeImageUsecase(url string, width, height int) (*model.ResizedImage, error) {
	s.logger.Info("Starting image resize")

	imageData, err := s.integration.GetImageApiReq(url)
	if err != nil {
		s.logger.Error("Failed to fetch image: " + err.Error())
		return nil, err
	}

	resizedImage, err := s.resizer.ResizeImage(imageData.Data, imageData.ContentType, width, height)
	if err != nil {
		s.logger.Error("Failed to resize image: " + err.Error())
		return nil, err
	}

	s.logger.Info("Image resized successfully")
	return resizedImage, nil
}
