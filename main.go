package main

import (
	"fmt"
	p "go-currency-exchange-rates/providers"
)

func main() {

	userInput := "Bank Of Canada"

	f, exists := p.MapProviders[userInput]
	if !exists {
		panic("provider not found")
	}
	fmt.Println(f())
}
