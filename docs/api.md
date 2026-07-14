# API

Planned endpoints:

- `GET /health` returns service health.
- `GET /` renders the upload page.
- `POST /upload` accepts a video upload and processing preset.
- `GET /media/{id}` returns media metadata or a result page.
- `GET /media/{id}/download` downloads the processed video.

The exact response formats should be finalized when the first vertical slice is implemented.
