package domain

import "time"

type TrackingTimeEntry struct {
	Id      TimeEntryId
	Minutes Minutes
	Since   time.Time
}

type StoppedTimeEntry struct {
	Id      TimeEntryId
	Minutes Minutes
}
