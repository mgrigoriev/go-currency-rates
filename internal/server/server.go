package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type Server struct {
	ratesCache *cache.Cache
	logger     *slog.Logger
	bindAddr   string
	tpl        *template.Template
}

func New(bindAddr string, ratesCache *cache.Cache, logger *slog.Logger) *Server {
	templatesDir := os.Getenv("TEMPLATES_DIR")
	if templatesDir == "" {
		templatesDir = "../templates"
	}

	return &Server{
		ratesCache: ratesCache,
		logger:     logger,
		bindAddr:   bindAddr,
		tpl:        template.Must(template.ParseFiles(templatesDir + "/index.html")),
	}
}

func (s *Server) ListenAndServe() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.indexHandler)
	r.HandleFunc("/from_rub/", s.conversionHandler)
	r.HandleFunc("/to_rub/", s.conversionHandler)
	http.Handle("/", r)

	s.logger.Info("Starting HTTP server at http://" + s.bindAddr)

	http.ListenAndServe(s.bindAddr, nil)
}

func (s *Server) indexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := s.tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		panic(err)
	}
}

func (s *Server) conversionHandler(w http.ResponseWriter, req *http.Request) {
	var result float64
	var fromCurrency string
	var toCurrency string

	amountParam := req.URL.Query().Get("amount")
	currencyParam := req.URL.Query().Get("currency")

	if amountParam == "" || currencyParam == "" {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
		return
	}

	amount, err := strconv.ParseFloat(amountParam, 64)
	if err != nil {
		http.Error(w, "Invalid amount value", http.StatusUnprocessableEntity)
		return
	}

	currency := strings.ToUpper(currencyParam)
	val, ok := s.ratesCache.Get(currency)
	if !ok {
		http.NotFound(w, req)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	if path == "from_rub" {
		result = amount / val
		fromCurrency = "RUB"
		toCurrency = currency
	} else {
		result = amount * val
		fromCurrency = currency
		toCurrency = "RUB"
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
