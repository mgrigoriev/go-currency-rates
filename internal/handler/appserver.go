package handler

import (
	"github.com/gorilla/mux"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"github.com/mgrigoriev/go-currency-rates/internal/converter"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"
)

type AppServer struct {
	converter *converter.Converter
	logger    *slog.Logger
	bindAddr  string
	tpl       *template.Template
}

func NewAppServer(bindAddr string, ratesCache *cache.Cache, logger *slog.Logger) *AppServer {
	templatesDir := os.Getenv("TEMPLATES_DIR")
	if templatesDir == "" {
		templatesDir = "../templates"
	}

	return &AppServer{
		converter: converter.NewConverter(ratesCache),
		logger:    logger,
		bindAddr:  bindAddr,
		tpl:       template.Must(template.ParseFiles(templatesDir + "/index.html")),
	}
}

func (s *AppServer) ListenAndServe() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.indexHandler).Methods("GET")
	r.HandleFunc("/convert/", s.convertHandler).Methods("GET")
	r.Use(s.loggingMiddleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         s.bindAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.logger.Info("Starting HTTP handler at http://" + s.bindAddr)

	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *AppServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("incoming request",
			"path", r.URL.Path,
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
		)

		next.ServeHTTP(w, r)
	})
}
