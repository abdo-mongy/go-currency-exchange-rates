package providers

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type ECBRate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

type ECBCubeTime struct {
	Time  string    `xml:"time,attr"`
	Rates []ECBRate `xml:"Cube"`
}

type ECBCubeRoot struct {
	TimeCubes []ECBCubeTime `xml:"Cube"`
}

type ECBEnvelope struct {
	XMLName xml.Name    `xml:"Envelope"`
	Cube    ECBCubeRoot `xml:"Cube"`
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

func (EuropeanCommercialBank) Prepare(data []byte) (currencies []Item) {
	var env ECBEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		panic(err)
	}
	timeCube := env.Cube.TimeCubes[0]

	for _, rate := range timeCube.Rates {
		date, _ := time.Parse("2006-01-02", timeCube.Time)
		currencies = append(currencies, Item{rate.Rate, date, rate.Currency})
	}
	return currencies
}

func (ecb EuropeanCommercialBank) GetRates() []Result {
	return Format(ecb.Prepare(ecb.Fetch()))
}

func init() {
	MapProviders["European Commercial Bank"] = EuropeanCommercialBank{}.GetRates
}
