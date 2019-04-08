package domain

import "time"

type TrackingTimeEntry struct {
	Id      string
	Minutes Minutes
	Since   time.Time
}

type StoppedTimeEntry struct {
	Id      string
	Minutes Minutes
}
