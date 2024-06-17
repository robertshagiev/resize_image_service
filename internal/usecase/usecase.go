package usecase

type ImageService struct {
	resizer Resizer
	logger  LoggerInterface
}

type Resizer interface {
	ResizeImage(url string, width, height int) ([]byte, string, error)
}

type LoggerInterface interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}

func NewImageService(resizer Resizer, log LoggerInterface) *ImageService {
	return &ImageService{
		resizer: resizer,
		logger:  log,
	}
}

func (s *ImageService) ResizeImage(url string, width, height int) ([]byte, string, error) {
	s.logger.Info("Starting image resize")
	data, format, err := s.resizer.ResizeImage(url, width, height)
	if err != nil {
		s.logger.Error("Failed to resize image: " + err.Error())
		return nil, "", err
	}
	s.logger.Info("Image resized successfully")
	return data, format, nil
}
