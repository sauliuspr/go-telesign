# go-telesign #

go-telesign is a Go client library for accessing the [Telesign API][https://developer.telesign.com/docs/api-docs]

go-telesign requires Go version 1.7 or greater

## Usage ##

```go
import "github.com/sauliuspr/go-telesign/telesign"
```

Construct a new Telesign client, the use the various service on the client to 
access different part of th Telesign API. For example:

```go

// send a SMS message
ts := telesign.APIAuthTransport{
    CustomerID: os.Getenv("TELESIGN_CUSTOMER_ID"),
    APIKey:     os.Getenv("TELESIGN_API_KEY"),
}
ctx := context.Background()
client := telesign.NewClient(ts.Client())

msg, _, err := client.Messaging.Send(ctx, "+14150001234", "Greetings Gophers!", "MKT", nil)

if err != nil {
    fmt.Println("ERROR: Failed to send messag via Messaging Telesign API: ", err)
} else {
    fmt.Println("Result:", msg)
}
```

### Authentication ###

The go-telesign has a separate package auth to handle authentication. 
