package mite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const userAgent = "mite-go/0.1 (+github.com/leanovate/mite-go)"
const layout = "2006-01-02"

type MiteApi interface {
	Projects() ([]*Project, error)
	Services() ([]Service, error)
	TimeEntries(params *TimeEntryParameters) ([]TimeEntry, error)
	TimeEntry(id string) (*TimeEntry, error)
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
