package main

import (
	"net/http"

	x402 "github.com/x402-foundation/x402/go"
	x402http "github.com/x402-foundation/x402/go/http"
	"github.com/x402-foundation/x402/go/http/nethttp"
	exactserver "github.com/x402-foundation/x402/go/mechanisms/evm/exact/server"
)

const (
	facilitator = "https://x402.org/facilitator"

	scheme  = "exact"
	payTo   = "0xYourAddress"
	price   = "0.01"
	network = "eip155:84532" // sepolia, base mainnet: "eip155:8453"
)

func main() {
	facilitator := x402http.NewFacilitatorClient(&x402http.FacilitatorConfig{
		URL: facilitator,
	})

	routes := x402http.RoutesConfig{
		"/resource": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme:  scheme,
					PayTo:   payTo,
					Price:   price,
					Network: x402.Network(network),
				},
			},
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/paid", nethttp.PaymentMiddlewareFromConfig(routes,
		nethttp.WithFacilitatorClient(facilitator),
		nethttp.WithScheme(x402.Network(network), exactserver.NewExactEvmScheme()),
	)(http.HandlerFunc(myHandler)))
	http.ListenAndServe(":8080", mux)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"payment accepted, here is the data"}`))
}
