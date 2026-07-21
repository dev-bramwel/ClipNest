package server

import (
	"html/template"
	"net/http"
	"path/filepath"

	"clipnest/internal/config"
	"clipnest/internal/handlers"
	"clipnest/internal/services"
)

func NewRouter(cfg *config.Config, templates *template.Template) http.Handler {
	uploadService := services.NewUploadService(filepath.Join("uploads"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.StaticDir))))
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/upload", handlers.NewUploadHandler(templates, uploadService))
	mux.HandleFunc("/result", handlers.NewResultHandler(templates, uploadService))
	mux.HandleFunc("/result/", handlers.NewResultHandler(templates, uploadService))
	mux.HandleFunc("/download/", handlers.NewDownloadHandler(uploadService))
	mux.HandleFunc("/thumbnail/", handlers.NewThumbnailHandler(uploadService))
	mux.HandleFunc("/", handlers.NewIndexHandler(templates))
	return mux
}
