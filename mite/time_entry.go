package mite

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type Direction int

const (
	DirectionAsc  = Direction(0)
	DirectionDesc = Direction(1)
)

type TimeEntryParameters struct {
	From      *time.Time
	To        *time.Time
	Direction *Direction
}

func (a *defaultApi) TimeEntries(params *TimeEntryParameters) ([]TimeEntry, error) {
	values := url.Values{}
	if params != nil {
		if params.From != nil {
			values.Add("from", params.From.Format(layout))
		}
		if params.To != nil {
			values.Add("to", params.To.Format(layout))
		}
		if params.Direction != nil {
			switch *params.Direction {
			case DirectionAsc:
				values.Add("direction", "asc")
			case DirectionDesc:
				values.Add("direction", "desc")
			}
		}
	}

	u, err := url.Parse(fmt.Sprintf("%s/%s", a.url, "time_entries.json"))
	if err != nil {
		return nil, err
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-MiteApiKey", a.key)
	req.Header.Add("User-Agent", userAgent)

	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	ter := []TimeEntryResponse{}
	err = json.NewDecoder(res.Body).Decode(&ter)
	if err != nil {
		return nil, err
	}

	timeEntries := []TimeEntry{}
	for _, te := range ter {
		timeEntries = append(timeEntries, te.ToTimeEntry())
	}

	return timeEntries, nil
}

type TimeEntryResponse struct {
	TimeEntry struct {
		Id          int    `json:"id"`
		Note        string `json:"note"`
		Minutes     int    `json:"minutes"`
		Date        string `json:"date_at"`
		ProjectName string `json:"project_name"`
		ServiceName string `json:"service_name"`
	} `json:"time_entry"`
}

func (r TimeEntryResponse) ToTimeEntry() TimeEntry {
	date, err := time.Parse(layout, r.TimeEntry.Date)
	if err != nil {
		panic(err)
	}

	return TimeEntry{
		Id:          fmt.Sprintf("%d", r.TimeEntry.Id),
		Note:        r.TimeEntry.Note,
		Duration:    time.Duration(r.TimeEntry.Minutes) * time.Minute,
		Date:        date,
		ProjectName: r.TimeEntry.ProjectName,
		ServiceName: r.TimeEntry.ServiceName,
	}
}
