package providers

import (
	"sort"
	"time"
)

type Result struct {
	Rate         float64
	Date         string
	CurrencyCode string
}

type Item struct {
	Rate         float64
	Date         time.Time
	CurrencyCode string
}

var MapProviders = make(map[string]func() []Result)

type Provider interface {
	Fetch() []byte
	Prepare([]byte) []Item
	GetRates() []Result
}

func Format(currencies []Item) (result []Result) {
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
