package main

import (
	"fmt"
	"github.com/mgrigoriev/go-currency-rates/internal/server"
	"github.com/mgrigoriev/go-currency-rates/internal/xmlparser"
)

const cbrApi = "https://www.cbr-xml-daily.ru/daily.xml"
const bindAddr = "0.0.0.0:9999"

func main() {
	prs := xmlparser.New(cbrApi)

	rates, err := prs.FetchAndStoreCurrencyRates()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	srv := server.New(bindAddr, rates)
	srv.ListenAndServe()
}
