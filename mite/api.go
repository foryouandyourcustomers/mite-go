package mite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/leanovate/mite-go/domain"
	"net/http"
	"net/url"
)

const contentType = "application/json"
const userAgentTemplate = "mite-go/%s (+github.com/leanovate/mite-go)"

type Api interface {
	domain.AccountApi
	domain.TimeEntryApi
	domain.TrackerApi
	domain.CustomerApi
	domain.ProjectApi
	domain.ServiceApi
	domain.UserApi
}

type api struct {
	base      string
	key       string
	userAgent string
	client    *http.Client
}

func NewApi(base string, key string, version string) (Api, error) {
	ua := fmt.Sprintf(userAgentTemplate, version)
	return &api{base: base, key: key, userAgent: ua, client: &http.Client{}}, nil
}

func (a *api) get(resource string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, a.url(resource), nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", a.userAgent)
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

func (a *api) getParametrized(resource string, values url.Values, result interface{}) error {
	u := &url.URL{}
	u.Path = resource
	u.RawQuery = values.Encode()

	return a.get(u.String(), result)
}

func (a *api) post(resource string, body interface{}, result interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.url(resource), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("User-Agent", a.userAgent)
	req.Header.Add("X-MiteApiKey", a.key)

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	if err := a.check(res); err != nil {
		return err
	}

	if result != nil {
		return json.NewDecoder(res.Body).Decode(result)
	}

	return nil
}

func (a *api) patch(resource string, body interface{}, result interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, a.url(resource), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("User-Agent", a.userAgent)
	req.Header.Add("X-MiteApiKey", a.key)

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	if err := a.check(res); err != nil {
		return err
	}

	if result != nil {
		return json.NewDecoder(res.Body).Decode(result)
	}

	return nil
}

func (a *api) delete(resource string, result interface{}) error {
	req, err := http.NewRequest(http.MethodDelete, a.url(resource), nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", a.userAgent)
	req.Header.Add("X-MiteApiKey", a.key)

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	if err := a.check(res); err != nil {
		return err
	}

	if result != nil {
		return json.NewDecoder(res.Body).Decode(result)
	}

	return nil
}

func (a *api) url(resource string) string {
	return fmt.Sprintf("%s/%s", a.base, resource)
}

func (a *api) check(res *http.Response) error {
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
