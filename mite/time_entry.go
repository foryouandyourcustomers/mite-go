package mite

import (
	"fmt"
	"github.com/leanovate/mite-go/date"
	"math"
	"net/url"
	"time"
)

type TimeEntry struct {
	Id          string
	Note        string
	Date        date.LocalDate
	Duration    time.Duration
	ProjectId   string
	ProjectName string
	ServiceId   string
	ServiceName string
}

type TimeEntryCommand struct {
	Date      *date.LocalDate
	Duration  *time.Duration
	Note      string
	ProjectId string
	ServiceId string
}

func (c *TimeEntryCommand) toRequest() *timeEntryRequest {
	r := &timeEntryRequest{}
	if c.Date != nil {
		r.TimeEntry.Date = c.Date.String()
	}
	if c.Duration != nil {
		r.TimeEntry.Minutes = int(math.Floor(math.Round(c.Duration.Minutes()))) // BOGUS
	}
	if c.Note != "" {
		r.TimeEntry.Note = c.Note
	}
	if c.ProjectId != "" {
		r.TimeEntry.ProjectId = c.ProjectId
	}
	if c.ServiceId != "" {
		r.TimeEntry.ServiceId = c.ServiceId
	}

	return r
}

type TimeEntryQuery struct {
	From      *date.LocalDate
	To        *date.LocalDate
	Direction string
}

func (q *TimeEntryQuery) toValues() url.Values {
	v := url.Values{}
	if q != nil {
		if q.From != nil {
			v.Add("from", q.From.String())
		}
		if q.To != nil {
			v.Add("to", q.To.String())
		}
		if q.Direction != "" {
			v.Add("direction", q.Direction)
		}
	}

	return v
}

type timeEntryRequest struct {
	TimeEntry struct {
		Date      string `json:"date_at"`
		Minutes   int    `json:"minutes"`
		Note      string `json:"note"`
		ProjectId string `json:"project_id"`
		ServiceId string `json:"service_id"`
	} `json:"time_entry"`
}

type timeEntryResponse struct {
	TimeEntry struct {
		Id          int    `json:"id"`
		Note        string `json:"note"`
		Minutes     int    `json:"minutes"`
		Date        string `json:"date_at"`
		ProjectId   int    `json:"project_id"`
		ProjectName string `json:"project_name"`
		ServiceId   int    `json:"service_id"`
		ServiceName string `json:"service_name"`
	} `json:"time_entry"`
}

func (r *timeEntryResponse) ToTimeEntry() *TimeEntry {
	d, err := date.Parse(r.TimeEntry.Date)
	if err != nil {
		panic(err)
	}

	return &TimeEntry{
		Id:          fmt.Sprintf("%d", r.TimeEntry.Id),
		Note:        r.TimeEntry.Note,
		Duration:    time.Duration(r.TimeEntry.Minutes) * time.Minute,
		Date:        d,
		ProjectId:   fmt.Sprintf("%d", r.TimeEntry.ProjectId),
		ProjectName: r.TimeEntry.ProjectName,
		ServiceId:   fmt.Sprintf("%d", r.TimeEntry.ServiceId),
		ServiceName: r.TimeEntry.ServiceName,
	}
}

func (a *api) TimeEntries(query *TimeEntryQuery) ([]*TimeEntry, error) {
	var ter []timeEntryResponse
	err := a.getParametrized("time_entries.json", query.toValues(), &ter)
	if err != nil {
		return nil, err
	}

	var timeEntries []*TimeEntry
	for _, te := range ter {
		timeEntries = append(timeEntries, te.ToTimeEntry())
	}

	return timeEntries, nil
}

func (a *api) TimeEntry(id string) (*TimeEntry, error) {
	ter := timeEntryResponse{}
	err := a.get(fmt.Sprintf("/time_entries/%s.json", id), &ter)
	if err != nil {
		return nil, err
	}

	return ter.ToTimeEntry(), nil
}

func (a *api) CreateTimeEntry(command *TimeEntryCommand) (*TimeEntry, error) {
	ter := timeEntryResponse{}
	err := a.post("/time_entries.json", command.toRequest(), &ter)
	if err != nil {
		return nil, err
	}

	return ter.ToTimeEntry(), nil
}

func (a *api) EditTimeEntry(id string, command *TimeEntryCommand) error {
	return a.patch(fmt.Sprintf("/time_entries/%s.json", id), command.toRequest())
}

func (a *api) DeleteTimeEntry(id string) error {
	return a.delete(fmt.Sprintf("/time_entries/%s.json", id))
}
