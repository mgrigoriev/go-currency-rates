package converter

import (
	"github.com/mgrigoriev/go-currency-rates/internal/testhelper"
	"testing"
)

func TestConvert(t *testing.T) {
	ratesCache := testhelper.InitCache(map[string]float64{"RUB": 1.0, "USD": 93.0, "EUR": 99.0})
	converter := NewConverter(ratesCache)

	tests := []struct {
		amount       float64
		fromCurrency string
		toCurrency   string
		expected     float64
		expectError  bool
	}{
		{100, "RUB", "EUR", 1.01, false},
		{200, "USD", "RUB", 18600.0, false},
		{100, "USD", "EUR", 93.94, false},
		{200, "JPY", "USD", 0.0, true},
		{150, "USD", "JPY", 0.0, true},
	}

	for _, test := range tests {
		result, err := converter.Convert(test.amount, test.fromCurrency, test.toCurrency)

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error, but got nil")
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		}
	}
}
