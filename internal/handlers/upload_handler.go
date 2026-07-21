package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"clipnest/internal/domain"
	"clipnest/internal/services"
)

type uploadPageData struct {
	ErrorMessage   string
	SuccessMessage string
	Preset         string
	FileName       string
	ResultURL      string
}

type resultPageData struct {
	Media         domain.Media
	Job           domain.Job
	ErrorMessage  string
	ThumbnailURL  string
	DownloadURL   string
	StatusText    string
	PresetLabel   string
	FileSizeText  string
	SuccessStatus bool
}

func NewUploadHandler(templates *template.Template, uploadService *services.UploadService) http.HandlerFunc {
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

			data.SuccessMessage = fmt.Sprintf("Uploaded %s and started %s processing.", result.OriginalName, presetLabel(string(result.Preset)))
			data.FileName = result.OriginalName
			if result.MediaID != "" {
				data.ResultURL = "/result/" + result.MediaID
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := templates.ExecuteTemplate(w, "layout.html", data); err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
			}
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func NewResultHandler(templates *template.Template, uploadService *services.UploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		mediaID := strings.Trim(strings.TrimPrefix(r.URL.Path, "/result"), "/")
		if mediaID == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := templates.ExecuteTemplate(w, "layout.html", resultPageData{ErrorMessage: "Missing media ID."}); err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
			}
			return
		}

		media, err := uploadService.GetMediaByID(mediaID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if executeErr := templates.ExecuteTemplate(w, "layout.html", resultPageData{ErrorMessage: err.Error()}); executeErr != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
			}
			return
		}

		job, _ := uploadService.GetJobByMediaID(mediaID)
		data := resultPageData{
			Media:         media,
			Job:           job,
			ThumbnailURL:  "/thumbnail/" + mediaID,
			DownloadURL:   "/download/" + mediaID,
			PresetLabel:   presetLabel(string(media.Preset)),
			FileSizeText:  formatBytes(media.SizeBytes),
			StatusText:    statusLabel(media.Status),
			SuccessStatus: media.Status == domain.MediaStatusReady,
		}
		if media.Status == domain.MediaStatusFailed {
			data.ErrorMessage = job.Error
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.ExecuteTemplate(w, "layout.html", data); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
	}
}

func NewDownloadHandler(uploadService *services.UploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		mediaID := strings.Trim(strings.TrimPrefix(r.URL.Path, "/download"), "/")
		path, err := uploadService.ResolveDownloadPath(mediaID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(path))
		http.ServeFile(w, r, path)
	}
}

func NewThumbnailHandler(uploadService *services.UploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		mediaID := strings.Trim(strings.TrimPrefix(r.URL.Path, "/thumbnail"), "/")
		media, err := uploadService.GetMediaByID(mediaID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if media.ThumbnailPath == "" {
			http.NotFound(w, r)
			return
		}
		if _, err := os.Stat(media.ThumbnailPath); err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, media.ThumbnailPath)
	}
}

func presetLabel(preset string) string {
	switch preset {
	case "reels":
		return "Reels"
	case "shorts":
		return "Shorts"
	case "whatsapp_status", "whatsapp-status", "whatsappstatus":
		return "WhatsApp Status"
	default:
		return "processing"
	}
}

func statusLabel(status domain.MediaStatus) string {
	switch status {
	case domain.MediaStatusReady:
		return "Ready"
	case domain.MediaStatusProcessing:
		return "Processing"
	case domain.MediaStatusFailed:
		return "Failed"
	default:
		return "Uploaded"
	}
}

func formatBytes(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
