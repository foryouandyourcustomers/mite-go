package mite

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct {
	Id   string
	Name string
	Note string
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
