package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	unit := flag.String("u", "toz", "The unit to convert the price to(g, kg)")
	currency := flag.String("c", "USD", "The currency to convert the price to(EUR, GBP)")

	flag.Parse()

	var symbol string
	switch strings.ToUpper(*currency) {
	case "USD":
		symbol = "$"
	case "EUR":
		symbol = "€"
	case "GBP":
		symbol = "£"
	default:
		fmt.Println("Invalid currency!")
		os.Exit(1)
	}
	goldPrice := getPrices(*unit, *currency)
	//fmt.Println(goldPrice)
	fmt.Printf("%s%s\n", strconv.FormatFloat(goldPrice, 'f', 2, 32), symbol)
}
func getPrices(unit string, currency string) float64 {
	resp, err := http.Get("https://api.gold-api.com/price/XAU")

	if err != nil {
		log.Fatal(err)
	}

	type priceResponse struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	var result priceResponse

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Println(readErr)
	}

	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}

	var price = float64(result.Price)
	//fmt.Println(price)

	switch unit {
	case "toz":
	case "g":
		price /= 31.103
	case "kg":
		price *= 32.15
	default:
		fmt.Println("Invalid unit!")
		os.Exit(1)
	}
	if strings.ToUpper(currency) != "USD" {
		exchangeRate := getExchangeRates(currency)
		price *= exchangeRate

	}

	defer resp.Body.Close()

	return price
}

func getExchangeRates(currency string) float64 {
	resp, err := http.Get("https://open.er-api.com/v6/latest/USD")

	if err != nil {
		log.Fatal(err)
	}

	type RatesResponse struct {
		Rates map[string]float64 `json:"rates"`
	}

	var result RatesResponse

	json.NewDecoder(resp.Body).Decode(&result)

	rate := result.Rates[strings.ToUpper(currency)]

	return float64(rate)
}

func printHelp() {
	fmt.Println("Usage: goldie <command>")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --u (name of unit)        Specify the unit to convert to if you dont like troy ounces")
	fmt.Println("Valid units:")
	fmt.Println("	g kg")
	fmt.Println("  --c (name of currency)       Specify the value to convert to if you dont like us dollars")
	fmt.Println(" Valid currencies:")
	fmt.Println("none right now lol")
}
