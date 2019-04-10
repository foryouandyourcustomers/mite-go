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
	base   *url.URL
	key    string
	agent  string
	client *http.Client
}

func NewApi(miteUrl string, miteKey string, clientVersion string) (Api, error) {
	base, err := url.Parse(miteUrl)
	if err != nil {
		return nil, err
	}

	userAgent := fmt.Sprintf(userAgentTemplate, clientVersion)

	return &api{base: base, key: miteKey, agent: userAgent, client: &http.Client{}}, nil
}

func (a *api) get(resource string, query url.Values, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, a.encode(resource, query), nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", a.agent)
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

func (a *api) post(resource string, body interface{}, result interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.encode(resource, nil), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("User-Agent", a.agent)
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

	req, err := http.NewRequest(http.MethodPatch, a.encode(resource, nil), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("User-Agent", a.agent)
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
	req, err := http.NewRequest(http.MethodDelete, a.encode(resource, nil), nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", a.agent)
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

func (a *api) encode(resource string, query url.Values) string {
	return a.base.ResolveReference(&url.URL{Path: resource, RawQuery: query.Encode()}).String()
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
