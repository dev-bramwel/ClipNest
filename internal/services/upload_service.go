package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"clipnest/internal/domain"
	"clipnest/internal/processor"
)

const defaultMaxUploadSize = 500 * 1024 * 1024

type UploadService struct {
	UploadDir   string
	MaxFileSize int64
	Processor   *processor.Processor
	MediaStore  map[string]domain.Media
	JobStore    map[string]domain.Job
}

type UploadResult struct {
	FilePath     string
	OutputPath   string
	OriginalName string
	Preset       domain.ProcessingPreset
	MediaID      string
	JobID        string
}

func NewUploadService(uploadDir string) *UploadService {
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	return &UploadService{
		UploadDir:   uploadDir,
		MaxFileSize: defaultMaxUploadSize,
		Processor:   processor.NewProcessor("ffmpeg"),
		MediaStore:  map[string]domain.Media{},
		JobStore:    map[string]domain.Job{},
	}
}

func (s *UploadService) HandleUpload(r *http.Request) (UploadResult, error) {
	if err := r.ParseMultipartForm(s.MaxFileSize); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			return UploadResult{}, errors.New("File exceeds the 500MB limit.")
		}
		return UploadResult{}, errors.New("Could not read the uploaded file.")
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) || strings.Contains(err.Error(), "no such file") {
			return UploadResult{}, errors.New("Please choose a video file.")
		}
		return UploadResult{}, errors.New("Please choose a video file.")
	}
	defer file.Close()

	presetValue := strings.TrimSpace(r.FormValue("preset"))
	presetConfig, err := processor.ResolvePreset(presetValue)
	if err != nil {
		return UploadResult{}, err
	}

	originalName := strings.TrimSpace(header.Filename)
	if originalName == "" {
		return UploadResult{}, errors.New("Please choose a video file.")
	}

	ext := strings.ToLower(filepath.Ext(originalName))
	if ext != ".mp4" && ext != ".mov" && ext != ".webm" {
		return UploadResult{}, errors.New("Please upload a supported video file: mp4, mov, or webm.")
	}

	if err := os.MkdirAll(s.UploadDir, 0o755); err != nil {
		return UploadResult{}, fmt.Errorf("could not prepare upload directory: %w", err)
	}
	if err := os.MkdirAll("processed", 0o755); err != nil {
		return UploadResult{}, fmt.Errorf("could not prepare processed directory: %w", err)
	}

	safeName := sanitizeFilename(originalName)
	baseName := strings.TrimSuffix(safeName, filepath.Ext(safeName))
	destPath := filepath.Join(s.UploadDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), safeName))
	destination, err := os.Create(destPath)
	if err != nil {
		return UploadResult{}, fmt.Errorf("could not save uploaded file: %w", err)
	}
	defer destination.Close()

	limitedReader := io.LimitReader(file, s.MaxFileSize+1)
	bytesWritten, err := io.Copy(destination, limitedReader)
	if err != nil {
		_ = os.Remove(destPath)
		return UploadResult{}, fmt.Errorf("could not save uploaded file: %w", err)
	}
	if bytesWritten > s.MaxFileSize {
		_ = os.Remove(destPath)
		return UploadResult{}, errors.New("File exceeds the 500MB limit.")
	}

	outputPath := filepath.Join("processed", fmt.Sprintf("%s_%s.mp4", baseName, strings.ToLower(string(presetConfig.Preset))))
	if s.Processor == nil {
		s.Processor = processor.NewProcessor("ffmpeg")
	}

	mediaID := fmt.Sprintf("media-%d", time.Now().UnixNano())
	jobID := fmt.Sprintf("job-%d", time.Now().UnixNano())
	media := domain.Media{
		ID:           mediaID,
		OriginalName: originalName,
		SourcePath:   destPath,
		OutputPath:   outputPath,
		Status:       domain.MediaStatusProcessing,
		Preset:       presetConfig.Preset,
		SizeBytes:    bytesWritten,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	job := domain.Job{
		ID:        jobID,
		MediaID:   mediaID,
		Type:      domain.JobTypeCompress,
		Status:    domain.JobStatusQueued,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.MediaStore[mediaID] = media
	s.JobStore[jobID] = job

	if err := s.Processor.Process(destPath, outputPath, presetConfig); err != nil {
		media.Status = domain.MediaStatusFailed
		media.UpdatedAt = time.Now()
		job.Status = domain.JobStatusFailed
		job.Error = err.Error()
		job.UpdatedAt = time.Now()
		s.MediaStore[mediaID] = media
		s.JobStore[jobID] = job
		return UploadResult{}, fmt.Errorf("processing failed: %s", err)
	}

	media.Status = domain.MediaStatusReady
	media.UpdatedAt = time.Now()
	job.Status = domain.JobStatusSucceeded
	job.UpdatedAt = time.Now()
	s.MediaStore[mediaID] = media
	s.JobStore[jobID] = job

	return UploadResult{FilePath: destPath, OutputPath: outputPath, OriginalName: originalName, Preset: presetConfig.Preset, MediaID: mediaID, JobID: jobID}, nil
}

func sanitizeFilename(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	base = strings.ReplaceAll(base, " ", "_")

	var builder strings.Builder
	for _, r := range base {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			builder.WriteRune(r)
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '.', r == '_', r == '-':
			builder.WriteRune(r)
		default:
			builder.WriteRune('_')
		}
	}

	cleaned := builder.String()
	if cleaned == "" {
		return "upload"
	}
	return cleaned
}
