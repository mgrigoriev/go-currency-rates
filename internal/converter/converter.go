package converter

import (
	"fmt"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
)

type Converter struct {
	ratesCache *cache.Cache
}

func NewConverter(ratesCache *cache.Cache) *Converter {
	return &Converter{ratesCache: ratesCache}
}

func (c *Converter) Convert(amount float64, fromCurrency string, toCurrency string) (float64, error) {
	fromRate, ok := c.ratesCache.Get(fromCurrency)
	if !ok {
		return 0, fmt.Errorf("currency rate not found for %s", fromCurrency)
	}
	toRate, ok := c.ratesCache.Get(toCurrency)
	if !ok {
		return 0, fmt.Errorf("currency rate not found for %s", toCurrency)
	}

	result := amount * fromRate / toRate

	return result, nil
}
