package mite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const userAgent = "mite-go/0.1 (+github.com/leanovate/mite-go)"
const layout = "2006-01-02"

type MiteApi interface {
	Projects() ([]Project, error)
	Services() ([]Service, error)
	TimeEntries(params *TimeEntryParameters) ([]TimeEntry, error)
}

type Project struct {
	Id   string
	Name string
	Note string
}

type Service struct {
	Id   string
	Name string
	Note string
}

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

type defaultApi struct {
	url    string
	key    string
	client *http.Client
}

func NewMiteApi(url string, key string) MiteApi {
	return &defaultApi{url: url, key: key, client: &http.Client{}}
}

func (a *defaultApi) Projects() ([]Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.url, "projects.json"), nil)
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

	prs := []ProjectResponse{}
	err = json.NewDecoder(res.Body).Decode(&prs)
	if err != nil {
		return nil, err
	}

	projects := []Project{}
	for _, pr := range prs {
		projects = append(projects, pr.ToProject())
	}

	return projects, nil
}

type ProjectResponse struct {
	Project struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Note string `json:"note"`
	} `json:"project"`
}

func (r ProjectResponse) ToProject() Project {
	return Project{
		Id:   fmt.Sprintf("%d", r.Project.Id),
		Name: r.Project.Name,
		Note: r.Project.Note,
	}
}

func (a *defaultApi) Services() ([]Service, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.url, "services.json"), nil)
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

	srs := []ServiceResponse{}
	err = json.NewDecoder(res.Body).Decode(&srs)
	if err != nil {
		return nil, err
	}

	services := []Service{}
	for _, sr := range srs {
		services = append(services, sr.ToService())
	}

	return services, nil
}

type ServiceResponse struct {
	Service struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Note string `json:"note"`
	} `json:"service"`
}

func (r ServiceResponse) ToService() Service {
	return Service{
		Id:   fmt.Sprintf("%d", r.Service.Id),
		Name: r.Service.Name,
		Note: r.Service.Note,
	}
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
