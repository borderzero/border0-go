package client

import (
	"encoding/json"
	"time"
)

type FlexibleTime struct {
	Time time.Time
}

func FlexibleTimeFrom(s string) (FlexibleTime, error) {
	t, err := time.Parse(time.RFC3339, s)
	return FlexibleTime{t}, err
}

func (f FlexibleTime) MarshalJSON() ([]byte, error) {
	if f.Time.IsZero() {
		return json.Marshal(0)
	}
	return json.Marshal(f.Time.Unix())
}

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

func (f FlexibleTime) String() string {
	if f.Time.IsZero() {
		return ""
	}
	return f.Time.Format(time.RFC3339)
}
