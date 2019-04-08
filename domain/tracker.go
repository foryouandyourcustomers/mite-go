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

type TrackerApi interface {
	Tracker() (*TrackingTimeEntry, error)
	StartTracker(id TimeEntryId) (*TrackingTimeEntry, *StoppedTimeEntry, error)
	StopTracker(id TimeEntryId) (*StoppedTimeEntry, error)
}
