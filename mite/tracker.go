package mite

import (
	"math"
	"strconv"
	"time"
)

type TrackingTimeEntry struct {
	Id       string
	Duration time.Duration
	Since    time.Time
}

type StoppedTimeEntry struct {
	Id       string
	Duration time.Duration
}

type Tracker struct {
	Tracking TrackingTimeEntry
	Stopped  StoppedTimeEntry
}

type TrackerCommand struct {
	Id       string
	Duration time.Duration
	Since    time.Time
}

func (c *TrackerCommand) toRequest() *trackerRequest {
	i, err := strconv.Atoi(c.Id)
	if err != nil {
		panic(err)
	}

	r := &trackerRequest{}
	r.Tracker.TrackingTimeEntry.Id = i
	r.Tracker.TrackingTimeEntry.Minutes = int(math.Floor(math.Round(c.Duration.Minutes()))) // BOGUS
	r.Tracker.TrackingTimeEntry.Since = c.Since

	return r
}

type trackerRequest struct {
	Tracker struct {
		TrackingTimeEntry struct {
			Id      int       `json:"id"`
			Minutes int       `json:"minutes"`
			Since   time.Time `json:"since"`
		} `json:"tracking_time_entry"`
	} `json:"tracker"`
}

func (a *api) Tracker() (*Tracker, error) {
	return &Tracker{}, nil
}

func (a *api) StartTracker(command TrackerCommand) (*Tracker, error) {
	return &Tracker{}, nil
}

func (a *api) StopTracker(id string) (*Tracker, error) {
	return &Tracker{}, nil
}
