# ClipNest

A Go web app that lets users upload a video, compress it, resize it for Reels/Shorts/WhatsApp Status, generate a thumbnail, and download the processed version.

Core application structs live in `internal/domain` so handlers, services, repositories, storage, and processors share one vocabulary.

## MVP features

- Upload a video file from the browser
- Pick a Reels, Shorts, or WhatsApp Status preset
- Validate file type and size before processing
- Save the original upload under uploads/
- Process and resize the video into a vertical export
- Save processed output under processed/
- Generate a thumbnail for the result
- Show processing status, metadata, and a download link

## Install FFmpeg

Install FFmpeg before running the app:

- Ubuntu/Debian: `sudo apt-get update && sudo apt-get install -y ffmpeg`
- macOS (Homebrew): `brew install ffmpeg`
- Windows: install from the official FFmpeg website and add it to your `PATH`

## Environment setup

Create a `.env` file in the project root based on `.env.example`:

```bash
cp .env.example .env
```

The app reads these values from the environment or `.env`:

- `HOST` (default `0.0.0.0`)
- `PORT` (default `8080`)
- `STATIC_DIR` (default `web/static`)
- `TEMPLATE_DIR` (default `web/templates`)

## Run the app

```bash
go run ./cmd/server
```

Then open `http://localhost:8080/`.

## Project structure

```text
clipnest/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   └── app.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── server/
│   │   ├── router.go
│   │   └── middleware.go
│   │
│   ├── handlers/
│   │   ├── health_handler.go
│   │   ├── media_handler.go
│   │   └── upload_handler.go
│   │
│   ├── services/
│   │   ├── job_service.go
│   │   ├── media_service.go
│   │   └── upload_service.go
│   │
│   ├── repositories/
│   │   ├── media_repository.go
│   │   └── user_repository.go
│   │
│   ├── domain/
│   │   └── domain.go
│   │
│   ├── storage/
│   │   ├── local_storage.go
│   │   └── storage.go
│   │
│   ├── processor/
│   │   ├── ffmpeg.go
│   │   └── processor.go
│   │
│   ├── database/
│   │   ├── db.go
│   │   └── migrations/
│   │
│   └── utils/
│       ├── response.go
│       ├── validator.go
│       └── filenames.go
│
├── web/
│   ├── static/
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   │
│   └── templates/
│       ├── layout.html
│       ├── index.html
│       ├── upload.html
│       └── result.html
│
├── uploads/
│   └── .gitkeep
│
├── processed/
│   └── .gitkeep
│
├── docs/
│   ├── architecture.md
│   └── api.md
│
├── tests/
│   └── integration/
│
├── .env.example
├── .gitignore
├── LICENSE
├── README.md
├── go.mod
└── go.sum
```
