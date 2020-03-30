package main

import "time"

func GetValidDateLowerAndUpperBound() (time.Time, time.Time) {
	now := time.Now()
	var lowerBound time.Time
	var upperBound time.Time
	if now.Hour() <= 12 {
		lowerBound = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		upperBound = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local)
	} else {
		lowerBound = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local)
		upperBound = time.Date(now.Year(), now.Month(), now.Day(), 24, 0, 0, 0, time.Local)
	}
	return lowerBound, upperBound
}
