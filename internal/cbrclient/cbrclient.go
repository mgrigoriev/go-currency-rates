package cbrclient

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/mgrigoriev/go-currency-rates/internal/cache"
	"golang.org/x/net/html/charset"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CBRClient struct {
	cbrURL     string
	ratesCache *cache.Cache
	logger     *slog.Logger
}

type CBRData struct {
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

func New(cbrURL string, ratesCache *cache.Cache, logger *slog.Logger) *CBRClient {
	return &CBRClient{
		cbrURL:     cbrURL,
		ratesCache: ratesCache,
		logger:     logger,
	}
}

func (c *CBRClient) FetchAndCacheRates() error {
	data, err := c.fetchRates(c.cbrURL)
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

	c.logger.Info("Fetched and cached currency rates")
	c.logger.Debug(fmt.Sprintf("%v", c.ratesCache))

	return nil
}

func (c *CBRClient) fetchRates(url string) ([]byte, error) {
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

func (c *CBRClient) unmarshalRates(data []byte) (*CBRData, error) {
	var cbrData CBRData

	xmlDecoder := xml.NewDecoder(bytes.NewReader(data))
	xmlDecoder.CharsetReader = charset.NewReaderLabel
	err := xmlDecoder.Decode(&cbrData)
	if err != nil {
		return nil, fmt.Errorf("error decoding XML: %v", err)
	}

	return &cbrData, nil
}

func (c *CBRClient) cacheRates(cbrData *CBRData) error {
	for _, entry := range cbrData.Entries {
		rate := c.normalizeRate(entry.VunitRate)
		c.ratesCache.Set(entry.CharCode, rate)
	}

	c.ratesCache.Set("RUB", 1.0)

	return nil
}

func (c *CBRClient) normalizeRate(rate string) float64 {
	val := strings.ReplaceAll(rate, ",", ".")
	res, _ := strconv.ParseFloat(val, 64)

	return res
}
