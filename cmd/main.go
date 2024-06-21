package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"resize_image_service/internal/config"
	"resize_image_service/internal/handler"
	"resize_image_service/internal/integration"
	"resize_image_service/internal/logger"
	"resize_image_service/internal/router"
	"resize_image_service/internal/service"
	"time"
)

func main() {
	log := logger.New()

	configPath, err := filepath.Abs("../config.yaml")
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to get absolute path: %v", err))
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load config: %v", err))
	}

	client := integration.NewHTTPClient(10 * time.Second)

	imageFetcher := integration.NewImageFetcher(client)

	imageService := service.NewImageService(imageFetcher, log)
	h := handler.NewHandler(imageService, log, cfg.MaxParallelRequests)

	r := router.NewRouter(h)

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info(fmt.Sprintf("Server is starting at %s", address))
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
