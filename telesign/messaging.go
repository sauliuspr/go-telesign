package telesign

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Send a message of message_type to the target phone_number
//
// See https://developer.telesign.com/docs/messaging-api for detailed API documentation.

type MessagingService service

type Messaging struct {
	ReferenceID     *string `json:"reference_id"`
	SubmitTimestamp *string `json:"submit_timestamp,omitempty"`
	Status          *Status `json:"status,omitempty"`
}

func (s *MessagingService) Send(ctx context.Context, phoneNumber string, message string, message_type string, v url.Values) (*Messaging, *http.Response, error) {

	rResp := new(Messaging)
	v = cleanValues(v)
	v.Set("phone_number", phoneNumber)
	v.Set("message", message)
	v.Set("message_type", message_type)

	data := v.Encode()
	req, err := s.client.NewRequest("POST", "/v1/messaging", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(v.Encode())))

	resp, err := s.client.Do(ctx, req, rResp)
	defer resp.Body.Close()
	if err != nil {
		return nil, resp, err
	}

	return rResp, resp, err
}

func (m Messaging) String() string {
	return Stringify(m)
}

func (s *MessagingService) Get(ctx context.Context, refID string) (*Messaging, *http.Response, error) {

	rResp := new(Messaging)

	req, err := s.client.NewRequest("GET", "/v1/messaging/"+refID, nil)

	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(ctx, req, rResp)
	defer resp.Body.Close()
	if err != nil {
		return nil, resp, err
	}

	return rResp, resp, err
}
