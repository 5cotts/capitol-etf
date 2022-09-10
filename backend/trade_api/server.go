package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// View the account info.
func getAccountInfo(w http.ResponseWriter, _ *http.Request) {
	var alpacaClient alpaca.Client = alpaca.NewClient(alpaca.ClientOpts{})

	acct, err := alpacaClient.GetAccount()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acct)
}

// See the trade positions of the current account.
func getPositions(w http.ResponseWriter, _ *http.Request) {
	var alpacaClient alpaca.Client = alpaca.NewClient(alpaca.ClientOpts{})

	pos, err := alpacaClient.ListPositions()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pos)
}

// Initiate the purchase of an asset.
func buy(w http.ResponseWriter, r *http.Request) {
	var alpacaClient alpaca.Client = alpaca.NewClient(alpaca.ClientOpts{})

	ticker := "AAPL"
	qty := decimal.NewFromFloat(1)
	resp, err := alpacaClient.PlaceOrder(
		alpaca.PlaceOrderRequest{
			AssetKey:    &ticker,
			Qty:         &qty,
			Side:        alpaca.Buy,
			Type:        alpaca.Market,
			TimeInForce: alpaca.GTC,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Initiate the sale of an asset.
func sell(_ http.ResponseWriter, _ *http.Request) {
	// Implement later
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/account_info", getAccountInfo).Methods("GET")
	router.HandleFunc("/positions", getPositions).Methods("GET")
	router.HandleFunc("/buy", buy).Methods("POST")
	router.HandleFunc("/sell", sell).Methods("POST")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))
}
