package main

import (
	"github.com/mgrigoriev/go-currency-rates/internal/cbrclient"
	"github.com/mgrigoriev/go-currency-rates/internal/server"
	"time"
)

const cbrApiUrl = "https://www.cbr-xml-daily.ru/daily.xml"
const bindAddr = "0.0.0.0:9999"

var ratesCache = make(map[string]float64)

func main() {
	client := cbrclient.New(cbrApiUrl, ratesCache)

	go func() {
		time.Sleep(5 * time.Second)
		err := client.FetchAndCacheRates()
		if err != nil {
			return
		}
	}()

	srv := server.New(bindAddr, ratesCache)
	srv.ListenAndServe()
}
