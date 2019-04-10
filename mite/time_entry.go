package mite

import (
	"fmt"
	"github.com/leanovate/mite-go/domain"
	"net/url"
	"time"
)

func fromCommand(c *domain.TimeEntryCommand) *timeEntryRequest {
	r := &timeEntryRequest{}
	if c.Date != nil {
		r.TimeEntry.Date = c.Date.String()
	}
	if c.Minutes != nil {
		r.TimeEntry.Minutes = c.Minutes.Value()
	}
	r.TimeEntry.Note = c.Note
	r.TimeEntry.UserId = int(c.UserId)
	r.TimeEntry.ProjectId = int(c.ProjectId)
	r.TimeEntry.ServiceId = int(c.ServiceId)
	r.TimeEntry.Locked = c.Locked

	return r
}

func fromQuery(q *domain.TimeEntryQuery) url.Values {
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
		Date      string `json:"date_at,omitempty"`
		Minutes   int    `json:"minutes,omitempty"`
		Note      string `json:"note,omitempty"`
		UserId    int    `json:"user_id,omitempty"`
		ProjectId int    `json:"project_id,omitempty"`
		ServiceId int    `json:"service_id,omitempty"`
		Locked    bool   `json:"locked,omitempty"`
	} `json:"time_entry"`
}

type timeEntryResponse struct {
	TimeEntry struct {
		Id           int       `json:"id"`
		Minutes      int       `json:"minutes"`
		Date         string    `json:"date_at"`
		Note         string    `json:"note"`
		Billable     bool      `json:"billable"`
		Locked       bool      `json:"locked"`
		Revenue      float64   `json:"revenue"`
		HourlyRate   int       `json:"hourly_rate"`
		UserId       int       `json:"user_id"`
		UserName     string    `json:"user_name"`
		ProjectId    int       `json:"project_id"`
		ProjectName  string    `json:"project_name"`
		CustomerId   int       `json:"customer_id"`
		CustomerName string    `json:"customer_name"`
		ServiceId    int       `json:"service_id"`
		ServiceName  string    `json:"service_name"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	} `json:"time_entry"`
}

func (r *timeEntryResponse) toTimeEntry() *domain.TimeEntry {
	d, err := domain.ParseLocalDate(r.TimeEntry.Date)
	if err != nil {
		panic(err)
	}

	return &domain.TimeEntry{
		Id:           domain.NewTimeEntryId(r.TimeEntry.Id),
		Minutes:      domain.NewMinutes(r.TimeEntry.Minutes),
		Date:         d,
		Note:         r.TimeEntry.Note,
		Billable:     r.TimeEntry.Billable,
		Locked:       r.TimeEntry.Locked,
		Revenue:      r.TimeEntry.Revenue,
		HourlyRate:   r.TimeEntry.HourlyRate,
		UserId:       domain.NewUserId(r.TimeEntry.UserId),
		UserName:     r.TimeEntry.UserName,
		ProjectId:    domain.NewProjectId(r.TimeEntry.ProjectId),
		ProjectName:  r.TimeEntry.ProjectName,
		CustomerId:   domain.NewCustomerId(r.TimeEntry.CustomerId),
		CustomerName: r.TimeEntry.CustomerName,
		ServiceId:    domain.NewServiceId(r.TimeEntry.ServiceId),
		ServiceName:  r.TimeEntry.ServiceName,
		CreatedAt:    r.TimeEntry.CreatedAt.UTC(),
		UpdatedAt:    r.TimeEntry.UpdatedAt.UTC(),
	}
}

func (a *api) TimeEntries(query *domain.TimeEntryQuery) ([]*domain.TimeEntry, error) {
	var ter []timeEntryResponse
	err := a.get("/time_entries.json", fromQuery(query), &ter)
	if err != nil {
		return nil, err
	}

	var timeEntries []*domain.TimeEntry
	for _, te := range ter {
		timeEntries = append(timeEntries, te.toTimeEntry())
	}

	return timeEntries, nil
}

func (a *api) TimeEntry(id domain.TimeEntryId) (*domain.TimeEntry, error) {
	ter := timeEntryResponse{}
	err := a.get(fmt.Sprintf("/time_entries/%s.json", id), nil, &ter)
	if err != nil {
		return nil, err
	}

	return ter.toTimeEntry(), nil
}

func (a *api) CreateTimeEntry(command *domain.TimeEntryCommand) (*domain.TimeEntry, error) {
	ter := timeEntryResponse{}
	err := a.post("/time_entries.json", fromCommand(command), &ter)
	if err != nil {
		return nil, err
	}

	return ter.toTimeEntry(), nil
}

func (a *api) EditTimeEntry(id domain.TimeEntryId, command *domain.TimeEntryCommand) error {
	return a.patch(fmt.Sprintf("/time_entries/%s.json", id), fromCommand(command), nil)
}

func (a *api) DeleteTimeEntry(id domain.TimeEntryId) error {
	return a.delete(fmt.Sprintf("/time_entries/%s.json", id), nil)
}
