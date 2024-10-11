package common

import "time"

var julianDaysOffset = 2440587.5
var secondsInADay = 86400.0

// JulianToTime converts a Julian date to time.Time
func JulianToTime(julianDate float64) time.Time {
	unixSeconds := (julianDate - julianDaysOffset) * secondsInADay
	return time.Unix(int64(unixSeconds), 0).UTC()
}

// TimeToJulian converts a time.Time object to Julian date
func TimeToJulian(t time.Time) float64 {
	unixSeconds := t.Unix()
	julianDate := float64(unixSeconds)/secondsInADay + julianDaysOffset
	return julianDate
}
