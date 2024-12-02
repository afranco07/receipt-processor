package receipt

import (
	"encoding/json"
	"time"
)

const timeOnly = "15:04"

// purchaseDate is custom type used to parse the receipt
// purchase date
type purchaseDate time.Time

func (d *purchaseDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}

	*d = purchaseDate(t)

	return nil
}

func (d *purchaseDate) scoreDay() int {
	// day is even
	if time.Time(*d).Day()%2 == 0 {
		return 0
	}

	return 6
}

// purchaseTime is custom type used to parse the receipt
// purchase time
type purchaseTime time.Time

func (pt *purchaseTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	t, err := time.Parse(timeOnly, s)
	if err != nil {
		return err
	}

	*pt = purchaseTime(t)

	return nil
}

func (pt *purchaseTime) scoreTime() int {
	hour := time.Time(*pt).Hour()

	// the time is between 2PM and 4PM
	if hour >= 14 && hour <= 16 {
		return 10
	}

	return 0
}
