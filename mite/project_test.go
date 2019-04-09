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

func TestApi_Projects(t *testing.T) {
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("[%s]", projectResponse)))
	}))

	defer srv.Close()

	api := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	projects, err := api.Projects()

	assert.Nil(t, err)
	assert.Equal(t, []*domain.Project{{
		Id:   domain.NewProjectId(643),
		Name: "Open-Source",
		Note: "valvat, memento et all."}}, projects)

	assert.Equal(t, http.MethodGet, rec.method)
	assert.Equal(t, "/projects.json", rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}
