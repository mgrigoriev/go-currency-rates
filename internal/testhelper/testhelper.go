package testhelper

import (
	"bytes"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"log/slog"
	"os"
)

func InitCache(rates map[string]float64) *cache.Cache {
	ratesCache := cache.New()

	for k, v := range rates {
		ratesCache.Set(k, v)
	}

	return ratesCache
}

func InitLogger() *slog.Logger {
	buffer := bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{}))

	return logger
}

func SetTemplatesDir(dir string) {
	err := os.Setenv("TEMPLATES_DIR", dir)
	if err != nil {
		panic(err)
	}
}
