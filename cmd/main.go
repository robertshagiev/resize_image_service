package main

import (
	"fmt"
	"net/http"
	"resize_image_service/internal/config"
	"resize_image_service/internal/handler"
	"resize_image_service/internal/logger"
	"resize_image_service/internal/resize"
	"resize_image_service/internal/router"
)

func main() {
	log := logger.New()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load config: %v", err))
	}

	resizer := resize.NewResizer()
	h := handler.NewHandler(resizer, log, cfg.MaxParallelRequests)

	r := router.NewRouter(h)

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info(fmt.Sprintf("Server is starting at %s", address))
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
