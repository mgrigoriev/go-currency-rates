package handler

import (
	"github.com/gorilla/mux"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"github.com/mgrigoriev/go-currency-rates/internal/converter"
	"log/slog"
	"net/http"
	"os"
	"text/template"
)

type HTTPServer struct {
	converter *converter.Converter
	logger    *slog.Logger
	bindAddr  string
	tpl       *template.Template
}

func NewHTTPServer(bindAddr string, ratesCache *cache.Cache, logger *slog.Logger) *HTTPServer {
	templatesDir := os.Getenv("TEMPLATES_DIR")
	if templatesDir == "" {
		templatesDir = "../templates"
	}

	return &HTTPServer{
		converter: converter.NewConverter(ratesCache),
		logger:    logger,
		bindAddr:  bindAddr,
		tpl:       template.Must(template.ParseFiles(templatesDir + "/index.html")),
	}
}

func (s *HTTPServer) ListenAndServe() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.indexHandler)
	r.HandleFunc("/convert/", s.convertHandler)
	http.Handle("/", r)

	s.logger.Info("Starting HTTP handler at http://" + s.bindAddr)

	http.ListenAndServe(s.bindAddr, nil)
}
