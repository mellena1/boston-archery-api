package db

import "time"

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
