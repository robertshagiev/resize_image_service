package router

import (
	"resize_image_service/internal/handler"

	"github.com/gorilla/mux"
)

func NewRouter(h *handler.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/resize", h.ResizeImage).Methods("GET")
	return r
}
