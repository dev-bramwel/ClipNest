package server

import (
	"html/template"
	"net/http"

	"clipnest/internal/config"
	"clipnest/internal/handlers"
)

func NewRouter(cfg *config.Config, templates *template.Template) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.StaticDir))))
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/upload", handlers.NewUploadHandler(templates))
	mux.HandleFunc("/", handlers.NewIndexHandler(templates))
	return mux
}
