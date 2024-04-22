package handler

import (
	"github.com/mgrigoriev/go-currency-rates/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListenAndServe(t *testing.T) {
	testhelper.SetTemplatesDir("../../templates")
	mockLogger := testhelper.InitLogger()
	mockCache := testhelper.InitCache(map[string]float64{})

	server := NewAppServer("localhost:8080", mockCache, mockLogger)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(server.indexHandler)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, recorder.Code)
	}
}
