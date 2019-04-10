package mite_test

import (
	"fmt"
	"github.com/leanovate/mite-go/domain"
	"github.com/leanovate/mite-go/mite"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const projectResponse = `{
  "project": {
    "id": 643,
    "name": "Open-Source",
    "note": "valvat, memento et all.",
    "customer_id": 291,
    "customer_name": "Yolk",
    "budget": 0,
    "budget_type": "minutes",
    "hourly_rate": 6000,
    "archived": false,
    "active_hourly_rate": "hourly_rate",
    "hourly_rates_per_service": [
      {
        "service_id": 31272,
        "hourly_rate": 4500
      },
      {
        "service_id": 149228,
        "hourly_rate": 5500
      }
    ],
    "created_at": "2011-08-17T12:06:57+02:00",
    "updated_at": "2015-02-19T10:53:10+01:00"
  }
}`

var projectObject = domain.Project{
	Id:   domain.NewProjectId(643),
	Name: "Open-Source",
	Note: "valvat, memento et all.",
}

func TestApi_Projects(t *testing.T) {
	// given
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(fmt.Sprintf("[%s]", projectResponse)).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	projects, err := api.Projects()

	// then
	assert.Nil(t, err)
	assert.Equal(t, []*domain.Project{&projectObject}, projects)

	assert.Equal(t, http.MethodGet, rec.RequestMethod())
	assert.Equal(t, "/projects.json", rec.RequestURI())
	assert.Empty(t, rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Empty(t, rec.RequestBody())
}
