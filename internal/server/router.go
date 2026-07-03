package server

func NewRouter(handlers *Handlers) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handlers.Health.Check)
	r.Get("/", handlers.Pages.Home)
	r.Post("/upload", handlers.Upload.Upload)

	return r
}
