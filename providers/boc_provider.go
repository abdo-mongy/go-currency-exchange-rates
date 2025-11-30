package providers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BOCData struct {
	XMLName xml.Name        `xml:"data"`
	O       BOCObservations `xml:"observations"`
}

type BOCObservations struct {
	Observations []BOCObservation `xml:"o"`
}

type BOCObservation struct {
	Date  string    `xml:"d,attr"`
	Rates []BOCRate `xml:"v"`
}

type BOCRate struct {
	Symbol string  `xml:"s,attr"`
	Rate   float64 `xml:",chardata"`
}

type BankOfCanada struct {
	Provider
}

func (BankOfCanada) Fetch() []byte {

	url := fmt.Sprintf("https://www.bankofcanada.ca/valet/observations/FXAUDCAD,FXBRLCAD,FXCNYCAD,FXEURCAD,FXHKDCAD,FXINRCAD,FXIDRCAD,FXJPYCAD,FXMYRCAD,FXMXNCAD,FXNZDCAD,FXNOKCAD,FXPENCAD,FXRUBCAD,FXSARCAD,FXSGDCAD,FXZARCAD,FXKRWCAD,FXSEKCAD,FXCHFCAD,FXTWDCAD,FXTHBCAD,FXTRYCAD,FXGBPCAD,FXUSDCAD,FXVNDCAD/xml?start_date=2017-01-03&end_date=%v", time.Now().Format("2006-01-02"))
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	return body
}

func (BankOfCanada) Prepare(data []byte) (currencies []Item) {
	var bocData BOCData
	if err := xml.Unmarshal(data, &bocData); err != nil {
		panic(err)
	}
	for _, o := range bocData.O.Observations {
		date, _ := time.Parse("2006-01-02", o.Date)
		for _, rate := range o.Rates {
			currencyCode := rate.Symbol[2 : len(rate.Symbol)-3]
			currencies = append(currencies, Item{rate.Rate, date, currencyCode})
		}
	}
	return currencies
}

func (boc BankOfCanada) GetRates() []Result {
	return Format(boc.Prepare(boc.Fetch()))
}

func init() {
	MapProviders["Bank Of Canada"] = BankOfCanada{}.GetRates
}
