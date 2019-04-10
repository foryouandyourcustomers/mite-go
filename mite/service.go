package mite

import (
	"github.com/leanovate/mite-go/domain"
)

type serviceResponse struct {
	Service struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Note string `json:"note"`
	} `json:"service"`
}

func (r *serviceResponse) toService() *domain.Service {
	return &domain.Service{
		Id:   domain.NewServiceId(r.Service.Id),
		Name: r.Service.Name,
		Note: r.Service.Note,
	}
}

func (a *api) Services() ([]*domain.Service, error) {
	var srs []serviceResponse
	err := a.get("/services.json", nil, &srs)
	if err != nil {
		return nil, err
	}

	var services []*domain.Service
	for _, sr := range srs {
		services = append(services, sr.toService())
	}

	return services, nil
}
