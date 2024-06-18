package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"resize_image_service/internal/api"
	"resize_image_service/internal/config"
	"resize_image_service/internal/handler"
	"resize_image_service/internal/logger"
	"resize_image_service/internal/resize"
	"resize_image_service/internal/router"
	"resize_image_service/internal/usecase"
)

func main() {
	log := logger.New()

	configPath, err := filepath.Abs("../internal/config/config.yaml")
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to get absolute path: %v", err))
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load config: %v", err))
	}

	apiCl := api.NewAPI()
	resizer := resize.NewResizer()
	imageService := usecase.NewImageService(apiCl, resizer, log)
	h := handler.NewHandler(imageService, log, cfg.MaxParallelRequests)

	r := router.NewRouter(h)

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info(fmt.Sprintf("Server is starting at %s", address))
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
