package mite

import (
	"fmt"
	"github.com/leanovate/mite-go/domain"
	"time"
)

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

func (r *trackerResponse) toTrackingTimeEntry() *domain.TrackingTimeEntry {
	if r.Tracker.TrackingTimeEntry == nil {
		return nil
	}

	return &domain.TrackingTimeEntry{
		Id:      domain.NewTimeEntryId(r.Tracker.TrackingTimeEntry.Id),
		Minutes: domain.NewMinutes(r.Tracker.TrackingTimeEntry.Minutes),
		Since:   r.Tracker.TrackingTimeEntry.Since.UTC(),
	}
}

func (r *trackerResponse) toStoppedTimeEntry() *domain.StoppedTimeEntry {
	if r.Tracker.StoppedTimeEntry == nil {
		return nil
	}

	return &domain.StoppedTimeEntry{
		Id:      domain.NewTimeEntryId(r.Tracker.StoppedTimeEntry.Id),
		Minutes: domain.NewMinutes(r.Tracker.StoppedTimeEntry.Minutes),
	}
}

func (a *api) Tracker() (*domain.TrackingTimeEntry, error) {
	tr := trackerResponse{}
	err := a.get("tracker.json", &tr)
	if err != nil {
		return nil, err
	}

	return tr.toTrackingTimeEntry(), nil
}

func (a *api) StartTracker(id domain.TimeEntryId) (*domain.TrackingTimeEntry, *domain.StoppedTimeEntry, error) {
	tr := &trackerResponse{}
	err := a.patch(fmt.Sprintf("tracker/%s.json", id), nil, tr)
	if err != nil {
		return nil, nil, err
	}

	return tr.toTrackingTimeEntry(), tr.toStoppedTimeEntry(), nil
}

func (a *api) StopTracker(id domain.TimeEntryId) (*domain.StoppedTimeEntry, error) {
	tr := &trackerResponse{}
	err := a.delete(fmt.Sprintf("tracker/%s.json", id), tr)
	if err != nil {
		return nil, err
	}

	return tr.toStoppedTimeEntry(), nil
}
