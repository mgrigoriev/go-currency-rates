package server

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	ratesCache map[string]float64
	bindAddr   string
}

func New(bindAddr string, ratesCache map[string]float64) *Server {
	return &Server{
		ratesCache: ratesCache,
		bindAddr:   bindAddr,
	}
}

func (s *Server) ListenAndServe() {
	http.HandleFunc("/", s.rates)
	http.ListenAndServe(s.bindAddr, nil)
}

func (s *Server) rates(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(s.ratesCache)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
