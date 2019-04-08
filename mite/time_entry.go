package mite

import (
	"fmt"
	"github.com/leanovate/mite-go/domain"
	"net/url"
	"strconv"
	"time"
)

type TimeEntry struct {
	Id           string
	Minutes      domain.Minutes
	Date         domain.LocalDate
	Note         string
	Billable     bool
	Locked       bool
	Revenue      float64
	HourlyRate   int
	UserId       string
	UserName     string
	ProjectId    string
	ProjectName  string
	CustomerId   string
	CustomerName string
	ServiceId    string
	ServiceName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TimeEntryCommand struct {
	Date      *domain.LocalDate
	Minutes   *domain.Minutes
	Note      string
	UserId    string
	ProjectId string
	ServiceId string
	Locked    bool
}

func (c *TimeEntryCommand) toRequest() *timeEntryRequest {
	r := &timeEntryRequest{}
	if c.Date != nil {
		r.TimeEntry.Date = c.Date.String()
	}
	if c.Minutes != nil {
		r.TimeEntry.Minutes = c.Minutes.Value()
	}
	r.TimeEntry.Note = c.Note
	r.TimeEntry.UserId = c.UserId
	r.TimeEntry.ProjectId = c.ProjectId
	r.TimeEntry.ServiceId = c.ServiceId
	r.TimeEntry.Locked = c.Locked

	return r
}

type TimeEntryQuery struct {
	From      *domain.LocalDate
	To        *domain.LocalDate
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
		Date      string `json:"date_at,omitempty"`
		Minutes   int    `json:"minutes,omitempty"`
		Note      string `json:"note,omitempty"`
		UserId    string `json:"user_id,omitempty"`
		ProjectId string `json:"project_id,omitempty"`
		ServiceId string `json:"service_id,omitempty"`
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

func (r *timeEntryResponse) toTimeEntry() *TimeEntry {
	d, err := domain.ParseLocalDate(r.TimeEntry.Date)
	if err != nil {
		panic(err)
	}

	return &TimeEntry{
		Id:           strconv.Itoa(r.TimeEntry.Id),
		Minutes:      domain.NewMinutes(r.TimeEntry.Minutes),
		Date:         d,
		Note:         r.TimeEntry.Note,
		Billable:     r.TimeEntry.Billable,
		Locked:       r.TimeEntry.Locked,
		Revenue:      r.TimeEntry.Revenue,
		HourlyRate:   r.TimeEntry.HourlyRate,
		UserId:       strconv.Itoa(r.TimeEntry.UserId),
		UserName:     r.TimeEntry.UserName,
		ProjectId:    strconv.Itoa(r.TimeEntry.ProjectId),
		ProjectName:  r.TimeEntry.ProjectName,
		CustomerId:   strconv.Itoa(r.TimeEntry.CustomerId),
		CustomerName: r.TimeEntry.CustomerName,
		ServiceId:    strconv.Itoa(r.TimeEntry.ServiceId),
		ServiceName:  r.TimeEntry.ServiceName,
		CreatedAt:    r.TimeEntry.CreatedAt,
		UpdatedAt:    r.TimeEntry.UpdatedAt,
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
		timeEntries = append(timeEntries, te.toTimeEntry())
	}

	return timeEntries, nil
}

func (a *api) TimeEntry(id string) (*TimeEntry, error) {
	ter := timeEntryResponse{}
	err := a.get(fmt.Sprintf("/time_entries/%s.json", id), &ter)
	if err != nil {
		return nil, err
	}

	return ter.toTimeEntry(), nil
}

func (a *api) CreateTimeEntry(command *TimeEntryCommand) (*TimeEntry, error) {
	ter := timeEntryResponse{}
	err := a.post("/time_entries.json", command.toRequest(), &ter)
	if err != nil {
		return nil, err
	}

	return ter.toTimeEntry(), nil
}

func (a *api) EditTimeEntry(id string, command *TimeEntryCommand) error {
	return a.patch(fmt.Sprintf("/time_entries/%s.json", id), command.toRequest(), nil)
}

func (a *api) DeleteTimeEntry(id string) error {
	return a.delete(fmt.Sprintf("/time_entries/%s.json", id), nil)
}
