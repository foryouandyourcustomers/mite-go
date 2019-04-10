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

const emptyTrackerResponse = `{"tracker":{}}`

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
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(trackingTimeEntryResponse).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	tracking, err := api.Tracker()

	// then
	assert.Nil(t, err)
	assert.Equal(t, &trackingTimeEntryObject, tracking)

	assert.Equal(t, http.MethodGet, rec.RequestMethod())
	assert.Equal(t, "/tracker.json", rec.RequestURI())
	assert.Empty(t, rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Empty(t, rec.RequestBody())
}

func TestApi_Tracker_Empty(t *testing.T) {
	// given
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(emptyTrackerResponse).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	tracking, err := api.Tracker()

	// then
	assert.Nil(t, err)
	assert.Nil(t, tracking)

	assert.Equal(t, http.MethodGet, rec.RequestMethod())
	assert.Equal(t, "/tracker.json", rec.RequestURI())
	assert.Empty(t, rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Empty(t, rec.RequestBody())
}

func TestApi_StartTracker(t *testing.T) {
	// given
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(trackingTimeEntryResponse).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	tracking, stopped, err := api.StartTracker(trackingTimeEntryObject.Id)

	// then
	assert.Nil(t, err)
	assert.Equal(t, &trackingTimeEntryObject, tracking)
	assert.Nil(t, stopped)

	assert.Equal(t, http.MethodPatch, rec.RequestMethod())
	assert.Equal(t, fmt.Sprintf("/tracker/%s.json", trackingTimeEntryObject.Id), rec.RequestURI())
	assert.Equal(t, "application/json", rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Equal(t, `{}`, rec.RequestBodyCanonical())
}

func TestApi_StartTracker_Running(t *testing.T) {
	// given
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(combinedTimeEntryResponse).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	tracking, stopped, err := api.StartTracker(trackingTimeEntryObject.Id)

	// then
	assert.Nil(t, err)
	assert.Equal(t, &trackingTimeEntryObject, tracking)
	assert.Equal(t, &stoppedTimeEntryObject, stopped)

	assert.Equal(t, http.MethodPatch, rec.RequestMethod())
	assert.Equal(t, fmt.Sprintf("/tracker/%s.json", trackingTimeEntryObject.Id), rec.RequestURI())
	assert.Equal(t, "application/json", rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Equal(t, `{}`, rec.RequestBodyCanonical())
}

func TestApi_StopTracker(t *testing.T) {
	// given
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(stoppedTimeEntryResponse).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	stopped, err := api.StopTracker(stoppedTimeEntryObject.Id)

	// then
	assert.Nil(t, err)
	assert.Equal(t, &stoppedTimeEntryObject, stopped)

	assert.Equal(t, http.MethodDelete, rec.RequestMethod())
	assert.Equal(t, fmt.Sprintf("/tracker/%s.json", stoppedTimeEntryObject.Id), rec.RequestURI())
	assert.Empty(t, rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Empty(t, rec.RequestBody())
}
