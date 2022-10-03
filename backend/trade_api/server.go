// A web-server written in Go to handle certain requests to the Alpaca API.
//
// Main website:
// 	https://alpaca.markets/
//
// API Documentation
// 	https://alpaca.markets/deprecated/docs/api-documentation/api-v2/
//
// Must have a `.env` file in root directory with the following variables.
// APCA_API_KEY_ID, APCA_API_SECRET_KEY, APCA_API_BASE_URL, PORT
//
// Execute the following to run the server.
// $ go run trade_api/server.go
//
// Sample request
// $ curl -X POST \
//	-d \ '{"symbol": "AAPL", "qty": 1}' \
// 	-H 'Content-Type: application/json' \
//  localhost:5000/buy

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

// Initiate the purchase of an asset with a POST request of the following shape:
// 	{"symbol": string, "qty": decimal}
func buy(w http.ResponseWriter, r *http.Request) {
	// Read the buy order from the POST request.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Marshall the POST request JSON into the correct struct.
	//
	// Instantiated with required fields according to the API documentation.
	// These defaults make sense for what we are building this for, so I don't think
	// they should be API parameters.
	//
	// See `[POST] Request a new order` below.
	// https://alpaca.markets/deprecated/docs/api-documentation/api-v2/orders/
	buy_order := alpaca.PlaceOrderRequest{
		Side:        alpaca.Buy,
		Type:        alpaca.Market,
		TimeInForce: alpaca.GTC,
	}
	if err := json.Unmarshal(body, &buy_order); err != nil {
		panic(err)
	}

	// Place the buy order
	var alpacaClient alpaca.Client = alpaca.NewClient(alpaca.ClientOpts{})
	resp, err := alpacaClient.PlaceOrder(buy_order)
	if err != nil {
		log.Fatal(err)
	}

	// Return the buy order response as JSON
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
