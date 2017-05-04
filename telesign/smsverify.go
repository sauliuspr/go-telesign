package telesign

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Sends a text message containing the verification code,
// to the specified phone number (supported for mobile phones only).
//
// Telesign API docs: https://developer.telesign.com/docs/rest_api-verify-sms

type SMSVerifyService service

type Status struct {
	UpdatedOn   *string `json:"updated_on"` // TODO convert to date
	Code        *int    `json:"code"`
	Description *string `json:"description"`
}
type Verify struct {
	CodeState   *string `json:"code_state"`
	CodeEntered *string `json:"code_entered"`
}
type PhoneType struct {
	Code        *int    `json:"code"`
	Description *string `json:"description"`
}
type Risk struct {
	Recommendation *string `json:"recommendation"`
	Score          *int    `json:"score"`
	Level          *string `json:"level"`
}
type Numbering struct {
	PhoneNumber *string `json:"phone_number"`
	MinLength   *int    `json:"min_length"`
	MaxLength   *int    `json:"max_length"`
	CountryCode *string `json:"country_code"`
}
type NumberDeactivationStatus struct {
	ErrorCode       int    `json:"error_code"`
	Description     string `json:"description"`
	LastDeactivated string `json:"last_deactivated"`
}
type SMSVerify struct {
	ReferenceID              *string                  `json:"reference_id"`
	ResourceURI              *string                  `json:"resource_uri"`
	SubResource              *string                  `json:"sub_resource"`
	Errors                   []Error                  `json:"errors"`
	Status                   *Status                  `json:"status,omitempty"`
	Verify                   *Verify                  `json:"verify,omitempty"`
	PhoneType                *PhoneType               `json:"phone_type,omitempty"`
	Risk                     *Risk                    `json:"risk,omitempty"`
	Numbering                *Numbering               `json:"numbering,omitempty"`
	NumberDeactivationStatus NumberDeactivationStatus `json:"number_deactivation_status,omitempty"`
}

type smsVerify SMSVerify

type smsVerifyUnmarshalHelper struct {
	smsVerify
	Attributes *smsVerify `json:"attributes"`
}

func (s *SMSVerify) UnmarshalJSON(b []byte) error {
	var helper smsVerifyUnmarshalHelper
	helper.Attributes = &helper.smsVerify
	if err := json.Unmarshal(b, &helper); err != nil {
		fmt.Println("ERROR: Unmarshaling error:", err)
		return err
	}
	*s = SMSVerify(helper.smsVerify)
	return nil
}

func (s *SMSVerifyService) Send(ctx context.Context, phoneNumber string, verifyCode string, language string, template string, numDeacReq bool, loopback bool) (*SMSVerify, *http.Response, error) {

	rResp := new(SMSVerify)

	form := url.Values{}
	form.Set("phone_number", phoneNumber)
	form.Add("language", language)
	form.Add("verify_code", verifyCode)
	form.Add("template", template)
	if loopback {
		form.Add("tst-provider", "RoutoSMPP_Loopback2")
	}
	if numDeacReq {
		form.Add("num_deac_req", "true")
	}
	data := form.Encode()
	req, err := s.client.NewRequest("POST", "/v1/verify/sms", bytes.NewBuffer([]byte(data)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
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
