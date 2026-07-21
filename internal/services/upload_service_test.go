package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleUploadRejectsMissingFile(t *testing.T) {
	service := NewUploadService(t.TempDir())

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("preset", "reels")
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, err := service.HandleUpload(req)
	if err == nil {
		t.Fatal("expected upload validation error")
	}
	if !strings.Contains(err.Error(), "Please choose a video file") {
		t.Fatalf("expected missing file message, got %v", err)
	}
}

func writeTestVideo(t *testing.T, dst io.Writer) error {
	t.Helper()

	videoPath := filepath.Join(t.TempDir(), "sample.mp4")
	cmd := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-f", "lavfi", "-i", "color=c=blue:s=320x240:d=1", "-frames:v", "1", "-f", "mp4", videoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	payload, err := os.ReadFile(videoPath)
	if err != nil {
		return err
	}
	if _, err := dst.Write(payload); err != nil {
		return err
	}
	return nil
}

func TestGetMediaByIDReturnsCleanErrorForMissingMedia(t *testing.T) {
	service := NewUploadService(t.TempDir())

	_, err := service.GetMediaByID("missing")
	if err == nil {
		t.Fatal("expected missing media error")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestResolveDownloadPathRejectsTraversal(t *testing.T) {
	service := NewUploadService(t.TempDir())

	_, err := service.ResolveDownloadPath("../bad")
	if err == nil {
		t.Fatal("expected traversal error")
	}
	if !strings.Contains(err.Error(), "invalid") {
		t.Fatalf("expected invalid path error, got %v", err)
	}
}

func TestHandleUploadAcceptsValidVideoFile(t *testing.T) {
	service := NewUploadService(t.TempDir())

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video", "sample.mp4")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if err := writeTestVideo(t, part); err != nil {
		t.Fatalf("write test video: %v", err)
	}
	_ = writer.WriteField("preset", "shorts")
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	result, err := service.HandleUpload(req)
	if err != nil {
		t.Fatalf("expected valid upload to succeed, got %v", err)
	}
	if result.OriginalName != "sample.mp4" {
		t.Fatalf("expected original name sample.mp4, got %q", result.OriginalName)
	}
	if _, err := os.Stat(result.FilePath); err != nil {
		t.Fatalf("expected uploaded file to exist: %v", err)
	}
	if _, err := os.Stat(result.OutputPath); err != nil {
		t.Fatalf("expected processed output to exist: %v", err)
	}
	if filepath.Ext(result.FilePath) != ".mp4" {
		t.Fatalf("expected mp4 extension, got %q", filepath.Ext(result.FilePath))
	}
}
