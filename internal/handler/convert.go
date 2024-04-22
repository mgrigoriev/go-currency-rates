package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *AppServer) convertHandler(w http.ResponseWriter, req *http.Request) {
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

	result, err = s.converter.Convert(amount, fromCurrency, toCurrency)
	if err != nil {
		http.Error(w, "Can't convert "+fromCurrency+" to "+toCurrency, http.StatusUnprocessableEntity)
	}

	response := map[string]interface{}{
		"amount":        amount,
		"from_currency": fromCurrency,
		"to_currency":   toCurrency,
		"result":        strconv.FormatFloat(result, 'f', 2, 64),
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
