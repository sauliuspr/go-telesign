package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/sauliuspr/go-telesign/telesign"
)

func main() {
	ts := telesign.APIAuthTransport{
		CustomerID: os.Getenv("TELESIGN_CUSTOMER_ID"),
		APIKey:     os.Getenv("TELESIGN_API_KEY"),
	}
	ctx := context.Background()
	client := telesign.NewClient(ts.Client())
	msg, _, err := client.Messaging.Send(ctx, "+37068600737", "Greetings Gophers!", "MKT", nil)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Messaging Reponse: %s\n", msg)
		// Wait 10 sec
		time.Sleep(10 * time.Second)
		sts, _, err := client.Messaging.Get(ctx, *msg.ReferenceID)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Messaging Status: %s\n", sts)
		}
	}
}
