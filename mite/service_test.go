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

func TestApi_Services(t *testing.T) {
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("[%s]", serviceResponse)))
	}))

	defer srv.Close()

	api := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	services, err := api.Services()

	assert.Nil(t, err)
	assert.Equal(t, []*domain.Service{{
		Id:   domain.NewServiceId(38672),
		Name: "Coding",
		Note: "will code for food"}}, services)

	assert.Equal(t, http.MethodGet, rec.method)
	assert.Equal(t, "/services.json", rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}
