package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"resize_image_service/internal/model"
	"strconv"
	"sync"
)

type Handler struct {
	service             imageService
	logger              logger
	mu                  sync.Mutex
	maxParallelRequests int
	currentRequests     int
}

type logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}

type imageService interface {
	ResizeImage(url string, width, height int) (*model.ResizedImage, error)
}

func NewHandler(service imageService, log logger, maxParallelRequests int) *Handler {
	return &Handler{
		service:             service,
		logger:              log,
		maxParallelRequests: maxParallelRequests,
		currentRequests:     0,
	}
}

func (h *Handler) ResizeImage(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	if h.currentRequests >= h.maxParallelRequests {
		h.mu.Unlock()
		http.Error(w, `{"error": "too many requests"}`, http.StatusTooManyRequests)
		return
	}
	h.currentRequests++
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		h.currentRequests--
		h.mu.Unlock()
	}()

	query := r.URL.Query()
	urlStr := query.Get("url")
	widthStr := query.Get("width")
	heightStr := query.Get("height")

	if urlStr == "" || widthStr == "" || heightStr == "" {
		h.resWithError(w, http.StatusBadRequest, "url, width and height are required parameters")
		return
	}

	decodedURL, err := url.QueryUnescape(urlStr)
	if err != nil {
		h.resWithError(w, http.StatusBadRequest, "invalid url parameter")
		return
	}

	width, err := strconv.Atoi(widthStr)
	if err != nil || width <= 0 {
		h.resWithError(w, http.StatusBadRequest, "width must be above 0")
		return
	}

	height, err := strconv.Atoi(heightStr)
	if err != nil || height <= 0 {
		h.resWithError(w, http.StatusBadRequest, "height must be above 0")
		return
	}

	resizedImage, err := h.service.ResizeImage(decodedURL, width, height)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to resize image: %v", err))
		h.resWithError(w, http.StatusInternalServerError, "failed to resize image")
		return
	}

	w.Header().Set("Content-Type", "image/"+resizedImage.Format)
	w.Write(resizedImage.Data)
}

func (h *Handler) resWithError(w http.ResponseWriter, code int, message string) {
	h.logger.Error(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
