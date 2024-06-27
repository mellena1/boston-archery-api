package handlers

import (
	"encoding/json"
	"time"
)

// Date is a type that will unmarhsal to a time.Time in the format of "2006-01-02"
// swagger:strfmt date
type Date time.Time

var _ json.Unmarshaler = &Date{}

func (d *Date) String() string {
	return d.ToTime().Format(time.DateOnly)
}

func (d *Date) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d Date) ToTime() time.Time {
	return time.Time(d)
}
