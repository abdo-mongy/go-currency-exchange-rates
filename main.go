package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"time"
)

type Result struct {
	Rate         float64
	Date         string
	CurrencyCode string
}

var PROVIDERS_MAP = map[string]func() []Result{
	"European Commercial Bank": EuropeanCommercialBank{}.GetRates,
}

var Providers []string = []string{"European Commercial Bank"}

type Provider struct{}

func (Provider) Fetch() []byte {
	return []byte{}
}

func (Provider) Format(data []byte) []Result {
	return []Result{}
}

func (p Provider) GetRates() []Result {
	return p.Format(p.Fetch())
}

type Item struct {
	Rate         float64
	Date         time.Time
	CurrencyCode string
}

type Rate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

type CubeTime struct {
	Time  string `xml:"time,attr"`
	Rates []Rate `xml:"Cube"`
}

type CubeRoot struct {
	TimeCubes []CubeTime `xml:"Cube"`
}

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Cube    CubeRoot `xml:"Cube"`
}

type EuropeanCommercialBank struct {
	Provider
}

func (EuropeanCommercialBank) Fetch() []byte {
	url := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	return body
}

func (EuropeanCommercialBank) Format(data []byte) []Result {

	var env Envelope
	if err := xml.Unmarshal(data, &env); err != nil {
		panic(err)
	}
	timeCube := env.Cube.TimeCubes[0]
	var currencies []Item
	var result []Result

	for _, rate := range timeCube.Rates {
		date, _ := time.Parse("2006-01-02", timeCube.Time)
		currencies = append(currencies, Item{rate.Rate, date, rate.Currency})
	}

	sort.Slice(currencies, func(i int, j int) bool {
		if currencies[i].Date != currencies[j].Date {
			return currencies[i].Date.Before(currencies[j].Date)
		}
		if currencies[i].CurrencyCode != currencies[j].CurrencyCode {
			return currencies[i].CurrencyCode < currencies[j].CurrencyCode
		}
		return currencies[i].Rate < currencies[j].Rate
	})

	for _, cur := range currencies {
		result = append(result, Result{cur.Rate, cur.Date.Format("2006-01-02"), cur.CurrencyCode})
	}

	return result
}

func (ecb EuropeanCommercialBank) GetRates() []Result {
	return ecb.Format(ecb.Fetch())
}

func main() {

	userInput := "European Commercial Bank"

	if !slices.Contains(Providers, userInput) {
		panic("Provider not found")
	}

	fmt.Println(PROVIDERS_MAP[userInput]())
}
