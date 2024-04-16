package xmlparser

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

type XmlParser struct {
	cbrApiUrl string
}

func New(cbrApiUrl string) *XmlParser {
	return &XmlParser{
		cbrApiUrl: cbrApiUrl,
	}
}

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

func (p *XmlParser) FetchAndStoreCurrencyRates() (map[string]float64, error) {
	cbrData, err := p.fetchCbrData(p.cbrApiUrl)
	if err != nil {
		return nil, err
	}

	rates, err := p.parseCbrData(cbrData)
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (p *XmlParser) fetchCbrData(url string) (*CbrCurrencyData, error) {
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

func (p *XmlParser) parseCbrData(cbrData *CbrCurrencyData) (map[string]float64, error) {
	rates := make(map[string]float64, len(cbrData.Entries))

	for _, entry := range cbrData.Entries {
		rateStr := strings.ReplaceAll(entry.VunitRate, ",", ".")
		rateFloat, err := strconv.ParseFloat(rateStr, 64)

		if err != nil {
			return nil, fmt.Errorf("error parsing currency rates: %v", err)
		}

		rates[entry.CharCode] = rateFloat
	}

	return rates, nil
}
