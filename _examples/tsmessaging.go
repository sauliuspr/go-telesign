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
	msg, _, err := client.Messaging.Send(ctx, "+37068600737", "test", "MKT", nil)

	if err != nil {
		fmt.Println("ERROR: Messaging Telesign API Error ", err)
	}
	fmt.Println("Result:", msg)
}
