package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type CbrCurrencyData struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Entries []struct {
		ID        string `xml:"ID,attr"`
		NumCode   string `xml:"NumCode"`
		CharCode  string `xml:"CharCode"`
		Nominal   string `xml:"Nominal"`
		Name      string `xml:"Name"`
		Value     string `xml:"Value"`
		VunitRate string `xml:"VunitRate"`
	} `xml:"Valute"`
}

func main() {
	cbrData, err := fetchCbrData("https://www.cbr-xml-daily.ru/daily.xml")
	if err != nil {
		fmt.Println("error fetching XML:", err)
		return
	}

	ratesCache := make(map[string]float64, len(cbrData.Entries))

	for _, entry := range cbrData.Entries {
		fmt.Printf("%v\n", entry)

		rateStr := strings.ReplaceAll(entry.VunitRate, ",", ".")
		rateFloat, err := strconv.ParseFloat(rateStr, 64)

		if err != nil {
			fmt.Println("error parsing currency rate:", err)
			return
		}

		ratesCache[entry.CharCode] = rateFloat
	}

	fmt.Println("USD Rate:", ratesCache["USD"])
	fmt.Println("EUR Rate:", ratesCache["EUR"])
}

func fetchCbrData(url string) (*CbrCurrencyData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	xmlDecoder := xml.NewDecoder(bytes.NewReader(data))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	var cbrData CbrCurrencyData
	err = xmlDecoder.Decode(&cbrData)
	if err != nil {
		return nil, fmt.Errorf("error decoding XML: %v", err)
	}

	return &cbrData, nil
}
