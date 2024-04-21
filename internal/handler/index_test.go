package handler

import (
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"log/slog"
	"net/http"
	"os"
	"testing"
)

func TestListenAndServe(t *testing.T) {
	os.Setenv("TEMPLATES_DIR", "../../templates")

	mockLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	mockCache := &cache.Cache{}

	server := NewHTTPServer("localhost:8080", mockCache, mockLogger)
	go server.ListenAndServe()

	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Errorf("Expected no error; got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, resp.StatusCode)
	}
}
