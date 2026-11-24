package jst

import "time"

const locationName = "Asia/Tokyo"

func Location() *time.Location {
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		return time.FixedZone("JST", 9*60*60)
	}
	return loc
}

func Now() time.Time {
	return time.Now().In(Location())
}

func In(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.In(Location())
}
