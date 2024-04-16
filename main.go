package main

import (
	"fmt"
)

func main() {
	var err error

	ratesCache, err = fetchCurrencyRates()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	listenAndServe()
}
