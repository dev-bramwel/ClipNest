package domain

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Media struct {
	ID            string
	OwnerID       string
	OriginalName  string
	SourcePath    string
	OutputPath    string
	ThumbnailPath string
	Status        MediaStatus
	Preset        ProcessingPreset
	SizeBytes     int64
	Duration      time.Duration
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type MediaStatus string

const (
	MediaStatusUploaded   MediaStatus = "uploaded"
	MediaStatusProcessing MediaStatus = "processing"
	MediaStatusReady      MediaStatus = "ready"
	MediaStatusFailed     MediaStatus = "failed"
)

type ProcessingPreset string

const (
	ProcessingPresetReels          ProcessingPreset = "reels"
	ProcessingPresetShorts         ProcessingPreset = "shorts"
	ProcessingPresetWhatsAppStatus ProcessingPreset = "whatsapp_status"
)

type Job struct {
	ID        string
	MediaID   string
	Type      JobType
	Status    JobStatus
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type JobType string

const (
	JobTypeCompress          JobType = "compress"
	JobTypeResize            JobType = "resize"
	JobTypeGenerateThumbnail JobType = "generate_thumbnail"
)

type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusSucceeded JobStatus = "succeeded"
	JobStatusFailed    JobStatus = "failed"
)
