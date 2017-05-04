package telesign

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Retrieves a score for the specified phone number.
// This ranks the phone number's "risk level" on a scale from 0 to 1000,
// so you can code your web application to handle particular use cases (e.g.,
// to stop things like chargebacks, identity theft, fraud, and spam).
//
// Telesign API docs: https://developer.telesign.com/docs/rest_api-phoneid-score

type PhoneIDScoreService service

type Location struct {
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       string `json:"zip"`
	MetroCode string `json:"metro_code"`
	County    string `json:"county"`
	Country   struct {
		Name string `json:"name"`
		Iso2 string `json:"iso2"`
		Iso3 string `json:"iso3"`
	} `json:"country"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	TimeZone struct {
		Name         string `json:"name"`
		UtcOffsetMin string `json:"utc_offset_min"`
		UtcOffsetMax string `json:"utc_offset_max"`
	} `json:"time_zone"`
}

type Carrier struct {
	Name *string `json:"name"`
}

type Fraud struct {
	FirstOccurredOn time.Time `json:"first_occurred_on"`
	FraudEvents     []struct {
		DiscoveredOn time.Time `json:"discovered_on"`
		FraudType    string    `json:"fraud_type"`
		Impact       string    `json:"impact"`
		ImpactType   string    `json:"impact_type"`
		Industry     string    `json:"industry"`
		OccurredOn   time.Time `json:"occurred_on"`
	} `json:"fraud_events"`
	LastOccurredOn        time.Time `json:"last_occurred_on"`
	MostFrequentFraudType string    `json:"most_frequent_fraud_type"`
	TotalIncidents        int       `json:"total_incidents"`
}

type PhoneIDScore struct {
	ReferenceID *string             `json:"reference_id"`
	ResourceURI *string             `json:"resource_uri"`
	SubResource *string             `json:"sub_resource"`
	Errors      []Error             `json:"errors"`
	Location    *Location           `json:"location,omitempty"`
	Status      *Status             `json:"status,omitempty"`
	Carrier     *Carrier            `json:"carrier,omitempty"`
	PhoneType   *PhoneType          `json:"phone_type,omitempty"`
	Risk        *Risk               `json:"risk,omitempty"`
	Numbering   *GetStatusNumbering `json:"numbering,omitempty"`
}

type phoneIDScore PhoneIDScore

type phoneIDScoreUnmarshalHelper struct {
	phoneIDScore
	Attributes *phoneIDScore `json:"attributes"`
}

func (s *PhoneIDScore) UnmarshalJSON(b []byte) error {
	var helper phoneIDScoreUnmarshalHelper
	helper.Attributes = &helper.phoneIDScore
	if err := json.Unmarshal(b, &helper); err != nil {
		fmt.Println("ERROR: Unmarshaling error:", err)
		return err
	}
	*s = PhoneIDScore(helper.phoneIDScore)
	return nil
}

func (s *PhoneIDScoreService) Get(ctx context.Context, phoneNumber string) (*PhoneIDScore, *http.Response, error) {

	rResp := new(PhoneIDScore)

	req, err := s.client.NewRequest("GET", "/v1/phoneid/score"+phoneNumber, nil)

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
