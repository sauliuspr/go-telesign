package telesign

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestPingService_Get(t *testing.T) {
	setup()
	defer teardown()
	ctx := context.Background()

	testMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{"data": "at"}`)
	})

	got, _, err := testClient.Ping.Get(ctx)
	if err != nil {
		t.Errorf("General.Ping returned error: %v", err)
	}

	expected := "PingedAt"
	if !strings.Contains(got.String(), expected) {
		t.Errorf("General.Ping returned %+v, want %+v", got, expected)
	}
}
