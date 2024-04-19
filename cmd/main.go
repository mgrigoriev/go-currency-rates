package main

import (
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"github.com/mgrigoriev/go-currency-rates/internal/cbrclient"
	"github.com/mgrigoriev/go-currency-rates/internal/server"
	"log/slog"
	"os"
	"time"
)

const cbrApiUrl = "https://www.cbr-xml-daily.ru/daily.xml"
const bindAddr = "0.0.0.0:9999"

func main() {
	logger := initLogger()
	ratesCache := cache.New()
	client := cbrclient.New(cbrApiUrl, ratesCache, logger)

	go func() {
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
	}()

	srv := server.New(bindAddr, ratesCache, logger)
	srv.ListenAndServe()
}

func initLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
