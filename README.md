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
client := telesign.NewClient(nil)

// send a SMS message

sms, _, err := client.Messaging.SMS(ctx, "+14150001234", nil)
```

### Authentication ###

The go-telesign has a package tsauth to handle authentication. 
