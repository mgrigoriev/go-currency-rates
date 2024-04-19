package main

import (
	"github.com/mgrigoriev/go-currency-rates/constants"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"github.com/mgrigoriev/go-currency-rates/internal/cbrclient"
	"github.com/mgrigoriev/go-currency-rates/internal/httpserver"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := initLogger()
	ratesCache := cache.New()
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
	srv := httpserver.New(constants.BindAddr, ratesCache, logger)
	srv.ListenAndServe()
}
