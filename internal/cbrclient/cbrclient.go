package cbrclient

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CbrClient struct {
	cbrApiUrl  string
	ratesCache map[string]float64
}

func New(cbrApiUrl string, ratesCache map[string]float64) *CbrClient {
	return &CbrClient{
		cbrApiUrl:  cbrApiUrl,
		ratesCache: ratesCache,
	}
}

type CbrData struct {
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

func (c *CbrClient) FetchAndCacheRates() error {
	data, err := c.fetchRates(c.cbrApiUrl)
	if err != nil {
		return err
	}

	cbrData, err := c.unmarshalRates(data)
	if err != nil {
		return err
	}

	err = c.cacheRates(cbrData)
	if err != nil {
		return err
	}

	fmt.Println(c.ratesCache)

	return nil
}

func (c *CbrClient) fetchRates(url string) ([]byte, error) {
	httpClient := http.Client{Timeout: 3 * time.Second}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return data, nil
}

func (c *CbrClient) unmarshalRates(data []byte) (*CbrData, error) {
	var cbrData CbrData

	xmlDecoder := xml.NewDecoder(bytes.NewReader(data))
	xmlDecoder.CharsetReader = charset.NewReaderLabel
	err := xmlDecoder.Decode(&cbrData)
	if err != nil {
		return nil, fmt.Errorf("error decoding XML: %v", err)
	}

	return &cbrData, nil
}

func (c *CbrClient) cacheRates(cbrData *CbrData) error {
	for _, entry := range cbrData.Entries {
		rate := c.normalizeRate(entry.VunitRate)
		c.ratesCache[entry.CharCode] = rate
	}

	return nil
}

func (c *CbrClient) normalizeRate(rate string) float64 {
	val := strings.ReplaceAll(rate, ",", ".")
	res, _ := strconv.ParseFloat(val, 64)

	return res
}
