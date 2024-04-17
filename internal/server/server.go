package server

import (
	"encoding/json"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	ratesCache *cache.Cache
	bindAddr   string
}

func New(bindAddr string, ratesCache *cache.Cache) *Server {
	return &Server{
		ratesCache: ratesCache,
		bindAddr:   bindAddr,
	}
}

func (s *Server) ListenAndServe() {
	http.HandleFunc("/", s.convert)
	http.ListenAndServe(s.bindAddr, nil)
}

func (s *Server) convert(w http.ResponseWriter, req *http.Request) {
	var roundedResult string

	w.Header().Set("Content-Type", "application/json")

	amountRubParam := req.URL.Query().Get("amount_rub")
	convertToParam := req.URL.Query().Get("convert_to")

	if amountRubParam == "" || convertToParam == "" {
		http.Error(w, "Usage example: /?amount_rub=1000&convert_to=USD", http.StatusUnprocessableEntity)
		return
	}

	amountRub, err := strconv.ParseFloat(amountRubParam, 64)
	if err != nil {
		http.Error(w, "Invalid amount_rub value", http.StatusUnprocessableEntity)
		return
	}

	convertTo := strings.ToUpper(convertToParam)
	val, ok := s.ratesCache.Get(convertTo)
	if !ok {
		http.NotFound(w, req)
		return
	}

	if floatValue, ok := val.(float64); ok {
		roundedResult = strconv.FormatFloat(amountRub/floatValue, 'f', 2, 64)
	} else {
		http.Error(w, "Invalid rate value", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"amount_rub": amountRub,
		"convert_to": convertTo,
		"result":     roundedResult,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
