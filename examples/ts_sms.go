package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
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
	sms, resp, err := client.SMSVerify.Send(ctx, "+37068600737", "123", "en", "NERK Code $$CODE$$", true, false)

	if err != nil {
		fmt.Println("ERROR: SMSVerify Telesign API Error ", err)
	} else {

		fmt.Println("SMS Verify Response:", resp.Status)
		fmt.Println("SMS Verify Status: ", *sms.Status.Code, *sms.Status.Description)
		fmt.Println("Reference ID ", *sms.ReferenceID)
		status, resp, err := client.GetStatus.Get(ctx, *sms.ReferenceID)
		if err != nil {
			fmt.Println("ERROR: GetStatus Telesign API Error ", err)

		} else {
			fmt.Println("Get Status Response:", resp.Status)
			fmt.Println("Get Status: ", *status.Status.Code, *status.Status.Description)
			if status.Verify != nil {
				fmt.Printf("Verify: %v\n", *status.Verify.CodeState)
			}
		}

		// Wait 5 sec
		time.Sleep(10 * time.Second)

		status, resp, err = client.GetStatus.Get(ctx, *sms.ReferenceID)
		if err != nil {
			fmt.Println("ERROR: GetStatus Telesign API Error ", err)

		} else {
			fmt.Println("Get Status Response:", resp.Status)
			fmt.Println("Get Status: ", *status.Status.Code, *status.Status.Description)
			if status.Verify != nil {
				fmt.Printf("Verify: %v\n\n", *status.Verify.CodeState)
			}
		}

		// Enter Code received
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Code: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		// Send Completion Data
		status, resp, err = client.GetStatus.Complete(ctx, *sms.ReferenceID, string(text))
		if err != nil {
			fmt.Println("ERROR: GetStatus Telesign API Error ", err)
			fmt.Println("Response Body:", resp)

		} else {
			fmt.Println("Get Status Response:", resp.Status)
			fmt.Println("Get Status: ", *status.Status.Code, *status.Status.Description)
			if status.Verify != nil {
				fmt.Printf("Verify: %v\n", *status.Verify.CodeState)
			}

		}
		// Wait 5 sec
		time.Sleep(20 * time.Second)

		status, resp, err = client.GetStatus.Get(ctx, *sms.ReferenceID)
		if err != nil {
			fmt.Println("ERROR: GetStatus Telesign API Error ", err)

		} else {
			fmt.Println("Get Status Response:", resp.Status)
			fmt.Println("Get Status: ", *status.Status.Code, *status.Status.Description)
			if status.Verify != nil {
				fmt.Printf("Verify: %v\n", *status.Verify.CodeState)
			}

		}
	}
}
