# Architecture

ClipNest is structured as a small layered Go web application.

- `cmd/server` contains the executable entrypoint.
- `internal/app` should compose configuration, storage, repositories, services, handlers, and the HTTP server.
- `internal/domain` contains core application data structs and enums shared across layers.
- `internal/handlers` should translate HTTP requests and responses.
- `internal/services` should own application workflows such as upload, processing, and job orchestration.
- `internal/repositories` should persist and retrieve domain records.
- `internal/storage` should read and write media files.
- `internal/processor` should wrap FFmpeg-specific video operations.
- `web/templates` and `web/static` contain server-rendered UI assets.

Domain types live in one package so the rest of the app can depend on the same vocabulary without duplicating structs between HTTP, persistence, and processing layers.
