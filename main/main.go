package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_ = alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    os.Getenv("AlpacaAPIKey"),
		APISecret: os.Getenv("AlpacaAPISecret"),
		BaseURL:   "https://paper-api.alpaca.markets/",
	})
	marketClient := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    os.Getenv("AlpacaAPIKey"),
		APISecret: os.Getenv("AlpacaAPISecret"),
	})

	fiveMinDur, _ := time.ParseDuration("10m")

	request := marketdata.GetCryptoBarsRequest{
		TimeFrame: marketdata.OneMin,
		Start:     time.Now().Add(-fiveMinDur),
		End:       time.Now(),
	}

	bars, err := marketClient.GetCryptoBars("DOGE/USD", request)
	if err != nil {
		panic(err)
	}
	for _, bar := range bars {
		fmt.Printf("%+v\n", bar)
	}
}
