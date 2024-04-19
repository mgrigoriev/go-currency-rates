package httpserver

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"github.com/mgrigoriev/go-currency-rates/internal/converter"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

type HTTPServer struct {
	ratesCache *cache.Cache
	logger     *slog.Logger
	bindAddr   string
	tpl        *template.Template
}

func New(bindAddr string, ratesCache *cache.Cache, logger *slog.Logger) *HTTPServer {
	templatesDir := os.Getenv("TEMPLATES_DIR")
	if templatesDir == "" {
		templatesDir = "../templates"
	}

	return &HTTPServer{
		ratesCache: ratesCache,
		logger:     logger,
		bindAddr:   bindAddr,
		tpl:        template.Must(template.ParseFiles(templatesDir + "/index.html")),
	}
}

func (s *HTTPServer) ListenAndServe() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.indexHandler)
	r.HandleFunc("/convert/", s.convertHandler)
	http.Handle("/", r)

	s.logger.Info("Starting HTTP httpserver at http://" + s.bindAddr)

	http.ListenAndServe(s.bindAddr, nil)
}

func (s *HTTPServer) indexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := s.tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		panic(err)
	}
}

func (s *HTTPServer) convertHandler(w http.ResponseWriter, req *http.Request) {
	var result float64

	amountParam := req.URL.Query().Get("amount")
	fromCurrency := req.URL.Query().Get("from")
	toCurrency := req.URL.Query().Get("to")

	if amountParam == "" || fromCurrency == "" || toCurrency == "" {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
		return
	}

	amount, err := strconv.ParseFloat(amountParam, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusUnprocessableEntity)
		return
	}

	cvt := converter.NewConverter(s.ratesCache)
	result, err = cvt.Convert(amount, fromCurrency, toCurrency)
	if err != nil {
		http.Error(w, "Can't convert "+fromCurrency+" to "+toCurrency, http.StatusUnprocessableEntity)
	}

	response := map[string]interface{}{
		"amount":        amount,
		"from_currency": fromCurrency,
		"to_currency":   toCurrency,
		"result":        strconv.FormatFloat(result, 'f', 2, 64),
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
