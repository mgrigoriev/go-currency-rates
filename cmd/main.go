package main

import (
	"github.com/mgrigoriev/go-currency-rates/constants"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"github.com/mgrigoriev/go-currency-rates/internal/cbrclient"
	"github.com/mgrigoriev/go-currency-rates/internal/handler"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := initLogger()
	ratesCache := initCache()
	client := cbrclient.New(constants.CBRURL, ratesCache, logger)

	go fetchAndCacheRates(client, logger)

	runHTTPServer(ratesCache, logger)
}

func initLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func initCache() *cache.Cache {
	ratesCache := cache.New()
	ratesCache.Set("RUB", 1.0)

	return ratesCache
}

func fetchAndCacheRates(client *cbrclient.CBRClient, logger *slog.Logger) {
	if err := client.FetchAndCacheRates(); err != nil {
		logger.Error(err.Error())
	}

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := client.FetchAndCacheRates(); err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

func runHTTPServer(ratesCache *cache.Cache, logger *slog.Logger) {
	srv := handler.NewHTTPServer(constants.BindAddr, ratesCache, logger)
	srv.ListenAndServe()
}
