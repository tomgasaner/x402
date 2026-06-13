package main

import (
	"net/http"

	x402 "github.com/x402-foundation/x402/go"
	x402http "github.com/x402-foundation/x402/go/http"
	"github.com/x402-foundation/x402/go/http/nethttp"
	exactserver "github.com/x402-foundation/x402/go/mechanisms/evm/exact/server"
)

func main() {
	facilitator := x402http.NewFacilitatorClient(&x402http.FacilitatorConfig{
		URL: "https://x402.org/facilitator",
	})

	routes := x402http.RoutesConfig{
		"/resource": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme:  "exact",
					PayTo:   "0xYourAddress",
					Price:   "0.01",
					Network: x402.Network("eip155:84532"), // base mainnet: x402.Network("eip155:8453"),
				},
			},
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/paid", nethttp.PaymentMiddlewareFromConfig(routes,
		nethttp.WithFacilitatorClient(facilitator),
		nethttp.WithScheme(x402.Network("eip155:*"), exactserver.NewExactEvmScheme()),
	)(http.HandlerFunc(myHandler)))
	http.ListenAndServe(":8080", mux)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"payment accepted, here is the data"}`))
}
