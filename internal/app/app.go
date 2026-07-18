package app

import (
	"context"
	"html/template"
	"net"
	"net/http"
	"path/filepath"

	"clipnest/internal/config"
	"clipnest/internal/server"
)

type App struct {
	config *config.Config
	router http.Handler
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	templates, err := parseTemplates(cfg.TemplateDir)
	if err != nil {
		return nil, err
	}

	return &App{
		config: cfg,
		router: server.NewRouter(cfg, templates),
	}, nil
}

func parseTemplates(dir string) (*template.Template, error) {
	return template.ParseGlob(filepath.Join(dir, "*.html"))
}

func (a *App) Run(ctx context.Context) error {
	addr := net.JoinHostPort(a.config.Host, a.config.Port)
	a.server = &http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	if ctx == nil {
		return a.server.ListenAndServe()
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return a.server.Shutdown(context.Background())
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
