package mite_test

import (
	"fmt"
	"github.com/leanovate/mite-go/domain"
	"github.com/leanovate/mite-go/mite"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const emptyResponse = `{"tracker":{}}`

const trackingTimeEntryResponse = `{
  "tracker": {
    "tracking_time_entry": {
      "id": 36135322,
      "minutes": 0,
      "since": "2015-10-15T17:33:52+02:00"
    }
  }
}`

const combinedTimeEntryResponse = `{
  "tracker": {
    "tracking_time_entry": {
      "id": 36135322,
      "minutes": 0,
      "since": "2015-10-15T17:33:52+02:00"
    },
    "stopped_time_entry": {
      "id": 36134329,
      "minutes": 46
    }
  }
}`

const stoppedTimeEntryResponse = `{
  "tracker": {
    "stopped_time_entry": {
      "id": 36134329,
      "minutes": 46
    }
  }
}`

var trackingTimeEntryObject = domain.TrackingTimeEntry{
	Id:      domain.NewTimeEntryId(36135322),
	Minutes: domain.NewMinutes(0),
	Since:   time.Date(2015, time.October, 15, 15, 33, 52, 0, time.UTC),
}

var stoppedTimeEntryObject = domain.StoppedTimeEntry{
	Id:      domain.NewTimeEntryId(36134329),
	Minutes: domain.NewMinutes(46),
}

func TestApi_Tracker(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(trackingTimeEntryResponse))
	}))

	defer srv.Close()

	api := mite.NewApi(srv.URL, testApiKey, testClientVersion)

	// when
	tracking, err := api.Tracker()

	// then
	assert.Nil(t, err)
	assert.Equal(t, &trackingTimeEntryObject, tracking)

	assert.Equal(t, http.MethodGet, rec.method)
	assert.Equal(t, "/tracker.json", rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}

func TestApi_Tracker_Empty(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(emptyResponse))
	}))

	defer srv.Close()

	api := mite.NewApi(srv.URL, testApiKey, testClientVersion)

	// when
	tracking, err := api.Tracker()

	// then
	assert.Nil(t, err)
	assert.Nil(t, tracking)

	assert.Equal(t, http.MethodGet, rec.method)
	assert.Equal(t, "/tracker.json", rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}

func TestApi_StartTracker(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(trackingTimeEntryResponse))
	}))

	defer srv.Close()

	api := mite.NewApi(srv.URL, testApiKey, testClientVersion)

	// when
	tracking, stopped, err := api.StartTracker(trackingTimeEntryObject.Id)

	// then
	assert.Nil(t, err)
	assert.Equal(t, &trackingTimeEntryObject, tracking)
	assert.Nil(t, stopped)

	assert.Equal(t, http.MethodPatch, rec.method)
	assert.Equal(t, fmt.Sprintf("/tracker/%s.json", trackingTimeEntryObject.Id), rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}

func TestApi_StartTracker_Running(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(combinedTimeEntryResponse))
	}))

	defer srv.Close()

	api := mite.NewApi(srv.URL, testApiKey, testClientVersion)

	// when
	tracking, stopped, err := api.StartTracker(trackingTimeEntryObject.Id)

	// then
	assert.Nil(t, err)
	assert.Equal(t, &trackingTimeEntryObject, tracking)
	assert.Equal(t, &stoppedTimeEntryObject, stopped)

	assert.Equal(t, http.MethodPatch, rec.method)
	assert.Equal(t, fmt.Sprintf("/tracker/%s.json", trackingTimeEntryObject.Id), rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}

func TestApi_StopTracker(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(stoppedTimeEntryResponse))
	}))

	defer srv.Close()

	api := mite.NewApi(srv.URL, testApiKey, testClientVersion)

	// when
	stopped, err := api.StopTracker(stoppedTimeEntryObject.Id)

	// then
	assert.Nil(t, err)
	assert.Equal(t, &stoppedTimeEntryObject, stopped)

	assert.Equal(t, http.MethodDelete, rec.method)
	assert.Equal(t, fmt.Sprintf("/tracker/%s.json", stoppedTimeEntryObject.Id), rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}
