package mite

import (
	"fmt"
)

type Service struct {
	Id   string
	Name string
	Note string
}

type serviceResponse struct {
	Service struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Note string `json:"note"`
	} `json:"service"`
}

func (r *serviceResponse) ToService() *Service {
	return &Service{
		Id:   fmt.Sprintf("%d", r.Service.Id),
		Name: r.Service.Name,
		Note: r.Service.Note,
	}
}

func (a *miteApi) Services() ([]*Service, error) {
	srs := []serviceResponse{}
	err := a.get("services.json", &srs)
	if err != nil {
		return nil, err
	}

	services := []*Service{}
	for _, sr := range srs {
		services = append(services, sr.ToService())
	}

	return services, nil
}
