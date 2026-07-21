package processor

import (
	"testing"

	"clipnest/internal/domain"
)

func TestResolvePresetReturnsVerticalDimensions(t *testing.T) {
	preset, err := ResolvePreset("reels")
	if err != nil {
		t.Fatalf("ResolvePreset returned error: %v", err)
	}
	if preset.Preset != domain.ProcessingPresetReels {
		t.Fatalf("expected reels preset, got %q", preset.Preset)
	}
	if preset.Width != 1080 || preset.Height != 1920 {
		t.Fatalf("expected 1080x1920 for reels, got %dx%d", preset.Width, preset.Height)
	}
}

func TestResolvePresetRejectsInvalidValue(t *testing.T) {
	if _, err := ResolvePreset("square"); err == nil {
		t.Fatal("expected invalid preset error")
	}
}
