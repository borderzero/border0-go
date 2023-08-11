package client

import (
	"encoding/json"
	"time"
)

// FlexibleTime is a time.Time that can be unmarshalled from either a string (RFC3339) or a number (unix timestamp).
// On marshalling, it is always marshalled as a number (unix timestamp). FlexibleTime is used for Border0 API connector
// token's `expires_at` and `created_at` fields.
type FlexibleTime struct {
	time.Time
}

// FlexibleTimeFrom returns a new FlexibleTime set to the given time from a string in RFC3339 format. It's a helper
// function for FlexibleTime.
func FlexibleTimeFrom(s string) (FlexibleTime, error) {
	t, err := time.Parse(time.RFC3339, s)
	return FlexibleTime{t}, err
}

// MarshalJSON marshals the FlexibleTime as a unix timestamp.
func (f FlexibleTime) MarshalJSON() ([]byte, error) {
	if f.IsZero() {
		return json.Marshal(0)
	}
	return json.Marshal(f.Unix())
}

// UnmarshalJSON unmarshals the FlexibleTime from either a string (RFC3339) or a number (unix timestamp).
func (f *FlexibleTime) UnmarshalJSON(data []byte) error {
	if data[0] == '"' { // string time: "2014-11-12T11:45:26.371Z"
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
		f.Time = t
	} else { // number time: 1415799926
		var i int64
		if err := json.Unmarshal(data, &i); err != nil {
			return err
		}
		f.Time = time.Unix(i, 0)
	}
	return nil
}

// String returns the FlexibleTime as a string in RFC3339 format. If the FlexibleTime is zero, it returns an empty string.
func (f FlexibleTime) String() string {
	if f.IsZero() {
		return ""
	}
	return f.Format(time.RFC3339)
}
