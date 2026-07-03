# ClipNest
A Go web app that lets users upload a video, compress it, resize it for Reels/Shorts/WhatsApp Status, generate a thumbnail, and download the processed version.
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
│   │   ├── upload_handler.go
│   │   └── media_handler.go
│   │
│   ├── services/
│   │   ├── upload_service.go
│   │   ├── media_service.go
│   │   └── job_service.go
│   │
│   ├── repositories/
│   │   ├── media_repository.go
│   │   └── user_repository.go
│   │
│   ├── models/
│   │   ├── media.go
│   │   ├── user.go
│   │   └── job.go
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