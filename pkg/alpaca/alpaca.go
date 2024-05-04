package alpaca

import (
	"fmt"
	"math"
	"os"
	"strings"

	alpaca "github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/shopspring/decimal"
)

type Alpaca struct {
	Client       *alpaca.Client
	MarketClient *marketdata.Client
	Tickers      Tickers
	Positions    map[string]int
}

type Tickers struct {
	TickerToName map[string]string
	NameToTicker map[string]string
}

func NewAlpaca() (*Alpaca, error) {
	client := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    os.Getenv("AlpacaAPIKey"),
		APISecret: os.Getenv("AlpacaAPISecret"),
		BaseURL:   "https://paper-api.alpaca.markets/",
	})
	marketClient := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    os.Getenv("AlpacaAPIKey"),
		APISecret: os.Getenv("AlpacaAPISecret"),
	})
	tickers, err := initTickers(client)
	if err != nil {
		return nil, fmt.Errorf("error initializing tickers")
	}

	return &Alpaca{
		Client:       client,
		MarketClient: marketClient,
		Tickers:      *tickers,
		Positions:    make(map[string]int),
	}, nil
}

func initTickers(client *alpaca.Client) (*Tickers, error) {
	assets, err := client.GetAssets(alpaca.GetAssetsRequest{
		Status: "active",
	})
	if err != nil {
		return nil, fmt.Errorf("error getting assets: %v", err)
	}
	tickerToName := make(map[string]string)
	nameToTicker := make(map[string]string)
	for _, asset := range assets {
		symb := strings.ToLower(asset.Symbol)
		name := strings.ToLower(asset.Name)
		tickerToName[symb] = name
		nameToTicker[name] = symb
	}
	return &Tickers{
		TickerToName: tickerToName,
		NameToTicker: nameToTicker,
	}, nil
}

func (a *Alpaca) AddPosition(security string, sentiment int) error {
	security = strings.ToLower(security)
	securityName := a.Tickers.TickerToName[security]
	if securityName == "" {
		return fmt.Errorf("security not found: %s", security)
	}
	if _, ok := a.Positions[security]; !ok {
		a.Positions[security] = sentiment
	} else {
		a.Positions[security] += sentiment
	}
	return nil
}

func (a *Alpaca) CloseAllPositions() error {
	_, err := a.Client.CloseAllPositions(alpaca.CloseAllPositionsRequest{
		CancelOrders: true,
	})
	return err
}

func (a *Alpaca) ExecutePositions() {
	for security, sentiment := range a.Positions {
		var err error
		var side string
		qty := decimal.NewFromInt(int64(math.Abs(float64(sentiment))))
		symbol := strings.ToUpper(security)
		if sentiment > 0 {
			_, err = a.Client.PlaceOrder(alpaca.PlaceOrderRequest{
				Symbol:      symbol,
				Side:        "buy",
				Qty:         &qty,
				Type:        "market",
				TimeInForce: "day",
			})
		}
		if sentiment < 0 {
			_, err = a.Client.PlaceOrder(alpaca.PlaceOrderRequest{
				Symbol:      symbol,
				Side:        "sell",
				Qty:         &qty,
				Type:        "market",
				TimeInForce: "day",
			})
		}
		fmt.Printf("Sent market order of | %s %s %s \n", qty.String(), symbol, side)
		if err != nil {
			fmt.Printf("Order for %s did not go through: %v\n", symbol, err)
		}
	}
	// reset all positions after execution
	a.Positions = make(map[string]int)
}
