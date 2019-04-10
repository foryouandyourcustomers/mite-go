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

const serviceResponse = `{
  "service": {
    "id": 38672,
    "name": "Coding",
    "note": "will code for food",
    "hourly_rate": 3300,
    "archived": false,
    "billable": true,
    "created_at": "2009-12-13T12:12:00+01:00",
    "updated_at": "2015-12-13T07:20:04+01:00"
  }
}`

var serviceObject = domain.Service{
	Id:   domain.NewServiceId(38672),
	Name: "Coding",
	Note: "will code for food",
}

func TestApi_Services(t *testing.T) {
	// given
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(fmt.Sprintf("[%s]", serviceResponse)).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	services, err := api.Services()

	// then
	assert.Nil(t, err)
	assert.Equal(t, []*domain.Service{&serviceObject}, services)

	assert.Equal(t, http.MethodGet, rec.RequestMethod())
	assert.Equal(t, "/services.json", rec.RequestURI())
	assert.Empty(t, rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Empty(t, rec.RequestBody())
}
