package mite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const contentType = "application/json"
const userAgent = "mite-go/0.1 (+github.com/leanovate/mite-go)"

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
	base   string
	key    string
	client *http.Client
}

func NewMiteApi(base string, key string) MiteApi {
	return &miteApi{base: base, key: key, client: &http.Client{}}
}

func (a *miteApi) get(resource string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, a.url(resource), nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-MiteApiKey", a.key)

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	if err := a.check(res); err != nil {
		return err
	}

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

	req, err := http.NewRequest(http.MethodPost, a.url(resource), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-MiteApiKey", a.key)

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	if err := a.check(res); err != nil {
		return err
	}

	return json.NewDecoder(res.Body).Decode(result)
}

func (a *miteApi) patch(resource string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, a.url(resource), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-MiteApiKey", a.key)

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	if err := a.check(res); err != nil {
		return err
	}

	return nil
}

func (a *miteApi) delete(resource string) error {
	req, err := http.NewRequest(http.MethodDelete, a.url(resource), nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-MiteApiKey", a.key)

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	if err := a.check(res); err != nil {
		return err
	}

	return nil
}

func (a *miteApi) url(resource string) string {
	return fmt.Sprintf("%s/%s", a.base, resource)
}

func (a *miteApi) check(res *http.Response) error {
	if res.StatusCode < 400 {
		return nil
	}

	msg := struct {
		Error string `json:"error"`
	}{}
	err := json.NewDecoder(res.Body).Decode(&msg)
	if err != nil {
		return fmt.Errorf("failed to %s %s", res.Request.Method, res.Request.RequestURI)
	}

	return fmt.Errorf("failed to %s %s: %s", res.Request.Method, res.Request.RequestURI, msg.Error)
}
