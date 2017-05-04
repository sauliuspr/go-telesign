package telesign

import (
	"time"
)

// Timestamp represents a time generated from a JSON string
type Timestamp struct {
	time.Time
}

// NewTimestamp creates a new Timestamp object from a ISO8601 date string
func NewTimestamp(date string) *Timestamp {
	ts, err := time.Parse(time.RFC3339, date)
	if err != nil {
		panic(err)
	}
	return &Timestamp{
		Time: ts,
	}
}

// String calls time.Time's String method
func (t Timestamp) String() string {
	return t.Time.String()
}

// UnmarshalJSON helps unmarshal ISO8601 dates in JSON
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var err error
	(*t).Time, err = time.Parse(`"`+time.RFC3339+`"`, string(data))
	return err
}
