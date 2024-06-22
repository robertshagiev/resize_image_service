package service

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"resize_image_service/internal/model"

	"github.com/disintegration/imaging"
)

type imageService struct {
	imageFetcher imageFetcher
	logger       logger
}

type imageFetcher interface {
	FetchImage(url string) (*model.ImageData, error)
}

type logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}

func NewImageService(imageFetcher imageFetcher, logger logger) *imageService {
	return &imageService{
		imageFetcher: imageFetcher,
		logger:       logger,
	}
}

func (s *imageService) ResizeImage(url string, width, height int) (*model.ResizedImage, error) {
	s.logger.Info("Starting image resize")

	imageData, err := s.imageFetcher.FetchImage(url)
	if err != nil {
		s.logger.Error("Failed to fetch image: " + err.Error())
		return nil, err
	}

	img, format, err := image.Decode(bytes.NewReader(imageData.Data))
	if err != nil {
		s.logger.Error("Failed to decode image: " + err.Error())
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
		s.logger.Error("Failed to encode resized image: " + err.Error())
		return nil, err
	}

	s.logger.Info("Image resized successfully")
	return &model.ResizedImage{
		Data:   buf.Bytes(),
		Format: format,
	}, nil
}
