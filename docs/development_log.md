# Development Log

This log records the main development decisions, implementation steps, and verification results for the ClipNest project. Each change should it be when a feature is introduced, change the architecture, adjust data models, or fix an important bug an entry should be made in as much clarity as possible.

## Entry Format

Use this structure for every new entry:

```md
## Day N - Short Title

**Date:** YYYY-MM-DD
**Author:** Name
**Branch:** branch-name

### Goal
Briefly explain what the work was meant to achieve.

### Implementation
- Describe the files or packages changed.
- Explain the important design decisions.
- Mention any constraints or assumptions.

### Verification
- List commands run, manual checks done, or tests added.
- Note anything that could not be tested.

### Next Steps
- List follow-up work that should happen after this entry.
```

## Day 1 - Domain Structure Cleanup

**Date:** 2026-07-15
**Author:** dev-bramwel
**Branch:** feat/datastructure

### Goal
Clean up the initial ClipNest scaffold and introduce a single domain module for the core application data structs. The aim was to make the project structure clearer before implementing the first working upload and processing flow.

### Implementation
- Added `internal/domain/domain.go` as the shared home for core structs and enums:
  - `User`
  - `Media`
  - `Job`
  - media statuses
  - job statuses
  - job types
  - processing presets
- Removed the empty `internal/models` files because the project will use `internal/domain` as the central vocabulary for core data.
- Normalized file names so each layer uses singular, consistent names:
  - `uploads_handler.go` became `upload_handler.go`
  - `uploads_service.go` became `upload_service.go`
  - `media_services.go` became `media_service.go`
  - `job_services.go` became `job_service.go`
  - `media_repositories.go` became `media_repository.go`
  - `uploads.html` became `upload.html`
- Added package declarations to empty Go files so each package is valid and the module can compile.
- Updated `.gitignore` so runtime media files in `uploads/` and `processed/` are ignored while keeping `.gitkeep` files tracked.
- Added starter environment values to `.env.example`.
- Added basic `Makefile` targets for `run`, `test`, `fmt`, and `vet`.
- Updated `README.md`, `docs/architecture.md`, and `docs/api.md` so the documentation reflects the current structure and planned MVP endpoints.

### Verification
- Ran `gofmt -w cmd internal` to format the Go files.
- Ran `go test ./...` successfully.
- No feature-level manual test was possible yet because the app is still a scaffold and the upload/processing flow has not been implemented.

### Next Steps
- Build the runnable HTTP server and `/health` endpoint.
- Implement the upload page and upload handler.
- Add video validation, safe file naming, and local storage.
- Integrate FFmpeg for compression, resizing presets, and thumbnail generation.
- Add the result page and processed video download endpoint.
