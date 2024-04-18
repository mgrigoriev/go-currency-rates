package main

import (
	"fmt"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"github.com/mgrigoriev/go-currency-rates/internal/cbrclient"
	"github.com/mgrigoriev/go-currency-rates/internal/server"
	"log/slog"
	"os"
)

const cbrApiUrl = "https://www.cbr-xml-daily.ru/daily.xml"
const bindAddr = "0.0.0.0:9999"

func main() {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	ratesCache := cache.New()
	client := cbrclient.New(cbrApiUrl, ratesCache)

	go func() {
		if err := client.FetchAndCacheRates(); err != nil {
			logger.Error(err.Error())
		}

		logger.Info(fmt.Sprintf("Fetched and cached rates: %v", ratesCache))
	}()

	logger.Info("Starting HTTP server at http://" + bindAddr)

	srv := server.New(bindAddr, ratesCache)
	srv.ListenAndServe()
}
