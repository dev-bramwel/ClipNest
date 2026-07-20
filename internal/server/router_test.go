package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"clipnest/internal/config"
)

func TestRouterServesIndexTemplate(t *testing.T) {
	templates := template.Must(template.ParseFiles(
		filepath.Join("..", "..", "web", "templates", "layout.html"),
		filepath.Join("..", "..", "web", "templates", "index.html"),
	))
	cfg := &config.Config{StaticDir: filepath.Join("..", "..", "web", "static"), TemplateDir: filepath.Join("..", "..", "web", "templates")}

	server := httptest.NewServer(NewRouter(cfg, templates))
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Fatalf("expected HTML content type, got %q", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	bodyText := string(body)
	if bodyText == "" {
		t.Fatalf("expected rendered HTML body")
	}
	if !strings.Contains(bodyText, "Make every clip feel polished.") {
		t.Fatalf("expected rendered home page content, got %q", bodyText)
	}
}

func TestRouterReturnsNotFoundForUnknownRoute(t *testing.T) {
	templates := template.Must(template.ParseFiles(
		filepath.Join("..", "..", "web", "templates", "layout.html"),
		filepath.Join("..", "..", "web", "templates", "index.html"),
	))
	cfg := &config.Config{StaticDir: filepath.Join("..", "..", "web", "static"), TemplateDir: filepath.Join("..", "..", "web", "templates")}

	server := httptest.NewServer(NewRouter(cfg, templates))
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/missing", server.URL))
	if err != nil {
		t.Fatalf("GET /missing failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}
