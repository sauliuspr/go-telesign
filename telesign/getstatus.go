package telesign

import (
	"context"
	"fmt"
	"net/http"
)

// GetStatusService Retrieves the verification result.
// You make this call in your web application after users complete the authentication transaction (using either a call or sms).
//
// Telesign API docs: https://developer.telesign.com/docs/rest_api-verify-transaction-callback
type GetStatusService service

type GetStatusNumbering struct {
	Original struct {
		PhoneNumber         *string `json:"phone_number"`
		CompletePhoneNumber *string `json:"complete_phone_number"`
		CountryCode         *string `json:"country_code"`
	} `json:"original"`
	Cleansing struct {
		Sms struct {
			CleansedCode *int    `json:"cleansed_code,omitempty"`
			CountryCode  *string `json:"country_code,omitempty"`
			MaxLength    *int    `json:"max_length,omitempty"`
			MinLength    *int    `json:"min_length,omitempty"`
			PhoneNumber  *string `json:"phone_number,omitempty"`
		} `json:"sms,omitempty"`
		Call struct {
			CleansedCode *int    `json:"cleansed_code,omitempty"`
			CountryCode  *string `json:"country_code,omitempty"`
			MaxLength    *int    `json:"max_length,omitempty"`
			MinLength    *int    `json:"min_length,omitempty"`
			PhoneNumber  *string `json:"phone_number,omitempty"`
		} `json:"call,omitempty"`
	} `json:"cleansing,omitempty"`
}

type GetStatus struct {
	ReferenceID *string             `json:"reference_id"`
	ResourceURI *string             `json:"resource_uri"`
	SubResource *string             `json:"sub_resource"`
	Errors      []Error             `json:"errors"`
	Status      *Status             `json:"status,omitempty"`
	Verify      *Verify             `json:"verify,omitempty"`
	PhoneType   *PhoneType          `json:"phone_type,omitempty"`
	Risk        *Risk               `json:"risk,omitempty"`
	Numbering   *GetStatusNumbering `json:"numbering,omitempty"`
}

type getStatus GetStatus

func (s *GetStatusService) Get(ctx context.Context, refID string) (*GetStatus, *http.Response, error) {

	rResp := new(GetStatus)

	req, err := s.client.NewRequest("GET", "/v1/verify/"+refID, nil)

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

func (s *GetStatusService) Complete(ctx context.Context, refID string, code string) (*GetStatus, *http.Response, error) {

	rResp := new(GetStatus)

	req, err := s.client.NewRequest("GET", "/v1/verify/"+refID+"?verify_code="+code, nil)

	fmt.Printf("Request: %v\n\n", req)

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
