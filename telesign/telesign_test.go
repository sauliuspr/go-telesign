package telesign

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sauliuspr/tsclient/telesign"
)

var (
	// testMux is the HTTP request multiplexer used with the test server.
	testMux *http.ServeMux

	// testClient is the cachet client being tested.
	testClient *Client

	// testServer is a test HTTP server used to provide mock API responses.
	testServer *httptest.Server
)

//type testValues map[string]string

// setup sets up a test HTTP server along with a telesign.Client that is configured to talk to that test server.
// Tests should register handlers on mux which provide mock responses for the API method being tested.
func setup() {

	ts := telesign.APIAuthTransport{
		CustomerID: os.Getenv("TELESIGN_CUSTOMER_ID"),
		APIKey:     os.Getenv("TELESIGN_API_KEY"),
	}

	if ts.CustomerID != "" && ts.APIKey != "" {
		return
	}

	// Test server
	testMux = http.NewServeMux()
	testServer = httptest.NewServer(testMux)

	// Cachet client configured to use test server
	testClient = NewClient(ts)
}

// teardown closes the test HTTP server.
func teardown() {
	testServer.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}
