package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"clipnest/internal/services"
)

type uploadPageData struct {
	ErrorMessage   string
	SuccessMessage string
	Preset         string
	FileName       string
}

func NewUploadHandler(templates *template.Template) http.HandlerFunc {
	uploadService := services.NewUploadService(filepath.Join("uploads"))

	return func(w http.ResponseWriter, r *http.Request) {
		data := uploadPageData{Preset: r.FormValue("preset")}

		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := templates.ExecuteTemplate(w, "layout.html", data); err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
			}
		case http.MethodPost:
			result, err := uploadService.HandleUpload(r)
			if err != nil {
				data.ErrorMessage = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				if executeErr := templates.ExecuteTemplate(w, "layout.html", data); executeErr != nil {
					http.Error(w, "template error", http.StatusInternalServerError)
				}
				return
			}

			data.SuccessMessage = fmt.Sprintf("Uploaded %s and queued it for %s processing.", result.OriginalName, presetLabel(string(result.Preset)))
			data.FileName = result.OriginalName
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := templates.ExecuteTemplate(w, "layout.html", data); err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
			}
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func presetLabel(preset string) string {
	switch preset {
	case "reels":
		return "Reels"
	case "shorts":
		return "Shorts"
	case "whatsapp_status":
		return "WhatsApp Status"
	default:
		return "processing"
	}
}
