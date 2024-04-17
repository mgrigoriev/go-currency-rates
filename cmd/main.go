package main

import (
	"fmt"
	"github.com/mgrigoriev/go-currency-rates/internal/cbrclient"
	"github.com/mgrigoriev/go-currency-rates/internal/server"
	"log/slog"
	"os"
	"time"
)

const cbrApiUrl = "https://www.cbr-xml-daily.ru/daily.xml"
const bindAddr = "0.0.0.0:9999"

var ratesCache = make(map[string]float64)

func main() {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	client := cbrclient.New(cbrApiUrl, ratesCache)

	go func() {
		time.Sleep(5 * time.Second)

		if err := client.FetchAndCacheRates(); err != nil {
			logger.Error(err.Error())
		}

		logger.Info(fmt.Sprintf("Cached rates: %#v", ratesCache))
	}()

	logger.Info("Starting HTTP server at http://" + bindAddr)

	srv := server.New(bindAddr, ratesCache)
	srv.ListenAndServe()
}
