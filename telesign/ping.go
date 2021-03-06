package telesign

import (
	"context"
	"net/http"
)

// PingService handles communication with the simple API service /ping
// It response with the data and signature string
type PingService service

// Representation of Ping service
type Ping struct {
	PingedAt  *string `json:"pinged_on"`
	Signature *string `json:"signature_string"`
	Errors    []Error `json:"errors"`
}

type ping Ping

func (p *PingService) Get(ctx context.Context) (*Ping, *http.Response, error) {
	req, err := p.client.NewRequest("GET", "ping", nil)
	if err != nil {
		return nil, nil, err
	}
	rResp := new(Ping)
	resp, err := p.client.Do(ctx, req, rResp)
	defer resp.Body.Close()
	if err != nil {
		return nil, resp, err
	}
	return rResp, resp, err
}

func (p Ping) String() string {
	return Stringify(p)
}
