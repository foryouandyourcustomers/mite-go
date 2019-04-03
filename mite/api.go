package mite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const userAgent = "mite-go/0.1 (+github.com/leanovate/mite-go)"
const layout = "2006-01-02"

type MiteApi interface {
	TimeEntries(query *TimeEntryQuery) ([]*TimeEntry, error)
	TimeEntry(id string) (*TimeEntry, error)
	CreateTimeEntry(command *TimeEntryCommand) (*TimeEntry, error)
	EditTimeEntry(id string, command *TimeEntryCommand) error
	DeleteTimeEntry(id string) error
	Projects() ([]*Project, error)
	Services() ([]*Service, error)
}

type miteApi struct {
	url    string
	key    string
	client *http.Client
}

func NewMiteApi(url string, key string) MiteApi {
	return &miteApi{url: url, key: key, client: &http.Client{}}
}

func (a *miteApi) get(resource string, result interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.url, resource), nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-MiteApiKey", a.key)
	req.Header.Add("User-Agent", userAgent)

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()

	return json.NewDecoder(res.Body).Decode(result)
}

func (a *miteApi) getParametrized(resource string, values url.Values, result interface{}) error {
	u := &url.URL{}
	u.Path = resource
	u.RawQuery = values.Encode()

	return a.get(u.String(), result)
}

func (a *miteApi) post(resource string, body interface{}, result interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", a.url, resource), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Add("X-MiteApiKey", a.key)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()

	return json.NewDecoder(res.Body).Decode(result)
}

func (a *miteApi) patch(resource string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", a.url, resource), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Add("X-MiteApiKey", a.key)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/json")

	res, err := a.client.Do(req)

	defer func() { _ = res.Body.Close() }()

	return err
}

func (a *miteApi) delete(resource string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", a.url, resource), nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-MiteApiKey", a.key)
	req.Header.Add("User-Agent", userAgent)

	res, err := a.client.Do(req)

	defer func() { _ = res.Body.Close() }()

	return err
}
