package mite

import (
	"fmt"
	"net/url"
	"time"
)

type TimeEntry struct {
	Id          string
	Note        string
	Duration    time.Duration
	Date        time.Time
	ProjectName string
	ServiceName string
}

type TimeEntryParameters struct {
	From      *time.Time
	To        *time.Time
	Direction string
}

type timeEntryResponse struct {
	TimeEntry struct {
		Id          int    `json:"id"`
		Note        string `json:"note"`
		Minutes     int    `json:"minutes"`
		Date        string `json:"date_at"`
		ProjectName string `json:"project_name"`
		ServiceName string `json:"service_name"`
	} `json:"time_entry"`
}

func (r *timeEntryResponse) ToTimeEntry() *TimeEntry {
	date, err := time.Parse(layout, r.TimeEntry.Date)
	if err != nil {
		panic(err)
	}

	return &TimeEntry{
		Id:          fmt.Sprintf("%d", r.TimeEntry.Id),
		Note:        r.TimeEntry.Note,
		Duration:    time.Duration(r.TimeEntry.Minutes) * time.Minute,
		Date:        date,
		ProjectName: r.TimeEntry.ProjectName,
		ServiceName: r.TimeEntry.ServiceName,
	}
}

func (a *miteApi) TimeEntries(params *TimeEntryParameters) ([]*TimeEntry, error) {
	values := url.Values{}
	if params != nil {
		if params.From != nil {
			values.Add("from", params.From.Format(layout))
		}
		if params.To != nil {
			values.Add("to", params.To.Format(layout))
		}
		if params.Direction != "" {
			values.Add("direction", params.Direction)
		}
	}

	ter := []timeEntryResponse{}
	err := a.getParametrized("time_entries.json", values, &ter)
	if err != nil {
		return nil, err
	}

	timeEntries := []*TimeEntry{}
	for _, te := range ter {
		timeEntries = append(timeEntries, te.ToTimeEntry())
	}

	return timeEntries, nil
}

func (a *miteApi) TimeEntry(id string) (*TimeEntry, error) {
	ter := timeEntryResponse{}
	err := a.get(fmt.Sprintf("/time_entries/%s.json", id), &ter)
	if err != nil {
		return nil, err
	}

	return ter.ToTimeEntry(), nil
}
