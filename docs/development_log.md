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

## Day 2 - Upload Flow and Processing UI

**Date:** 2026-07-21
**Author:** dev-bramwel
**Branch:** feat/datastructure

### Goal
Turn the initial scaffold into a working MVP flow for receiving uploads, validating them, storing them locally, and presenting a result page for processed media.

### Implementation
- Added upload, result, download, and thumbnail handlers in the web layer and wired them through the main router.
- Implemented the upload service with multipart form parsing, file size validation, safe filename sanitization, local storage under `uploads/`, and processed output generation under `processed/`.
- Added FFmpeg-backed processing support for compression presets, output naming, and thumbnail creation.
- Introduced in-memory media and job tracking so the app can follow each upload through queued, running, ready, and failed states.
- Expanded the templates and frontend assets to support the upload form, processing feedback, and result page links.

### Verification
- Ran `go test ./...` successfully.
- Reviewed the router and handler wiring for the new upload and result flow.

### Next Steps
- Persist media and job data beyond the current in-memory store.
- Move processing work out of the request path into a background job queue.
- Add end-to-end tests for the upload and result experience.
