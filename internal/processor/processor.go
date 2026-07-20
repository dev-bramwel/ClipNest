package processor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"clipnest/internal/domain"
)

type PresetConfig struct {
	Preset     domain.ProcessingPreset
	Label      string
	Width      int
	Height     int
	VideoCodec string
	Bitrate    string
}

type Processor struct {
	FFMPEGPath string
}

func NewProcessor(ffmpegPath string) *Processor {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	return &Processor{FFMPEGPath: ffmpegPath}
}

func ResolvePreset(value string) (PresetConfig, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "reels":
		return PresetConfig{Preset: domain.ProcessingPresetReels, Label: "Reels", Width: 1080, Height: 1920, VideoCodec: "libx264", Bitrate: "2500k"}, nil
	case "shorts":
		return PresetConfig{Preset: domain.ProcessingPresetShorts, Label: "Shorts", Width: 1080, Height: 1920, VideoCodec: "libx264", Bitrate: "2200k"}, nil
	case "whatsapp-status", "whatsapp_status", "whatsappstatus":
		return PresetConfig{Preset: domain.ProcessingPresetWhatsAppStatus, Label: "WhatsApp Status", Width: 720, Height: 1280, VideoCodec: "libx264", Bitrate: "1800k"}, nil
	default:
		return PresetConfig{}, fmt.Errorf("Please choose a valid preset: Reels, Shorts, or WhatsApp Status.")
	}
}

func (p *Processor) Process(inputPath, outputPath string, preset PresetConfig) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("could not create processed directory: %w", err)
	}

	if _, err := exec.LookPath(p.FFMPEGPath); err != nil {
		return copyFile(inputPath, outputPath)
	}

	args := []string{
		"-y",
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2", preset.Width, preset.Height, preset.Width, preset.Height),
		"-c:v", preset.VideoCodec,
		"-crf", "23",
		"-preset", "medium",
		"-b:v", preset.Bitrate,
		"-c:a", "aac",
		"-b:a", "128k",
		outputPath,
	}

	cmd := exec.Command(p.FFMPEGPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("processing failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return nil
}

func copyFile(inputPath, outputPath string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("could not open input file: %w", err)
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer output.Close()

	if _, err := io.Copy(output, input); err != nil {
		return fmt.Errorf("could not copy processed file: %w", err)
	}

	return nil
}
