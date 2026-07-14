package config

type Config struct {
	Port            string
	UploadDir       string
	ProcessedDir    string
	MaxUploadSizeMB int
	DatabaseURL     string
	FFmpegPath      string
}
