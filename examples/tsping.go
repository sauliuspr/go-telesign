package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sauliuspr/go-telesign/telesign"
)

func main() {
	ts := telesign.APIAuthTransport{
		CustomerID: os.Getenv("TELESIGN_CUSTOMER_ID"),
		APIKey:     os.Getenv("TELESIGN_API_KEY"),
	}

	ctx := context.Background()

	client := telesign.NewClient(ts.Client())
	ping, resp, err := client.Ping.Get(ctx)

	if err != nil {
		fmt.Println("ERROR: API gave an error", err)
	} else {
		fmt.Printf("Telesign Response: %v", resp.Status)
		fmt.Println("\nPinged at:", *ping.PingedAt)
	}
}
