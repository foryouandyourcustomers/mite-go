package mite

import (
	"fmt"
	"github.com/leanovate/mite-go/datetime"
	"strconv"
	"time"
)

type TrackingTimeEntry struct {
	Id      string
	Minutes datetime.Minutes
	Since   time.Time
}

type StoppedTimeEntry struct {
	Id      string
	Minutes datetime.Minutes
}

type trackerResponse struct {
	Tracker struct {
		TrackingTimeEntry *struct {
			Id      int       `json:"id"`
			Minutes int       `json:"minutes"`
			Since   time.Time `json:"since"`
		} `json:"tracking_time_entry"`
		StoppedTimeEntry *struct {
			Id      int `json:"id"`
			Minutes int `json:"minutes"`
		} `json:"stopped_time_entry"`
	} `json:"tracker"`
}

func (r *trackerResponse) toTrackingTimeEntry() *TrackingTimeEntry {
	if r.Tracker.TrackingTimeEntry == nil {
		return nil
	}

	return &TrackingTimeEntry{
		Id:      strconv.Itoa(r.Tracker.TrackingTimeEntry.Id),
		Minutes: datetime.NewMinutes(r.Tracker.TrackingTimeEntry.Minutes),
		Since:   r.Tracker.TrackingTimeEntry.Since,
	}
}

func (r *trackerResponse) toStoppedTimeEntry() *StoppedTimeEntry {
	if r.Tracker.StoppedTimeEntry == nil {
		return nil
	}

	return &StoppedTimeEntry{
		Id:      strconv.Itoa(r.Tracker.StoppedTimeEntry.Id),
		Minutes: datetime.NewMinutes(r.Tracker.StoppedTimeEntry.Minutes),
	}
}

func (a *api) Tracker() (*TrackingTimeEntry, error) {
	tr := trackerResponse{}
	err := a.get("/tracker.json", &tr)
	if err != nil {
		return nil, err
	}

	return tr.toTrackingTimeEntry(), nil
}

func (a *api) StartTracker(id string) (*TrackingTimeEntry, *StoppedTimeEntry, error) {
	tr := &trackerResponse{}
	err := a.patch(fmt.Sprintf("/tracker/%s.json", id), nil, tr)
	if err != nil {
		return nil, nil, err
	}

	return tr.toTrackingTimeEntry(), tr.toStoppedTimeEntry(), nil
}

func (a *api) StopTracker(id string) (*StoppedTimeEntry, error) {
	tr := &trackerResponse{}
	err := a.delete(fmt.Sprintf("/tracker/%s.json", id), tr)
	if err != nil {
		return nil, err
	}

	return tr.toStoppedTimeEntry(), nil
}
