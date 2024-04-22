package handler

import (
	"github.com/mgrigoriev/go-currency-rates/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestConvertHandler(t *testing.T) {
	testhelper.SetTemplatesDir("../../templates")
	mockLogger := testhelper.InitLogger()
	mockCache := testhelper.InitCache(map[string]float64{"RUB": 1.0, "USD": 93.0, "EUR": 99.0})

	server := NewAppServer("localhost:8080", mockCache, mockLogger)

	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
	}{
		{
			name: "ValidConversion",
			queryParams: map[string]string{
				"amount": "100",
				"from":   "USD",
				"to":     "EUR",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "InvalidAmount",
			queryParams: map[string]string{
				"amount": "invalid",
				"from":   "USD",
				"to":     "EUR",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "MissingParams",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusMovedPermanently,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := url.Values{}
			for key, value := range tt.queryParams {
				query.Set(key, value)
			}

			req, err := http.NewRequest("GET", "/convert?"+query.Encode(), nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(server.convertHandler)
			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d; got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}
