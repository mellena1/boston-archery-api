package db

import (
	"time"

	"github.com/mellena1/boston-archery-api/slices"
)

func dateToString(t time.Time) string {
	return t.Format(time.DateOnly)
}

func stringToDate(s string) time.Time {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		// should never have an invalid date string in the DB
		panic(err)
	}
	return t
}

func dateSliceToStrs(dates []time.Time) []string {
	return slices.Map(dates, func(date time.Time) string {
		return dateToString(date)
	})
}

func strSliceToDates(dates []string) []time.Time {
	return slices.Map(dates, func(s string) time.Time {
		return stringToDate(s)
	})
}
