package mite_test

import (
	"fmt"
	"github.com/leanovate/mite-go/domain"
	"github.com/leanovate/mite-go/mite"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const timeEntryResponse = `{
  "time_entry": {
    "id": 36159117,
    "minutes": 15,
    "date_at": "2015-10-16",
    "note": "Rework description of authentication process",
    "billable": true,
    "locked": false,
    "revenue": null,
    "hourly_rate": 0,
    "user_id": 211,
    "user_name": "Noah Scott",
    "project_id": 88309,
    "project_name": "API Docs",
    "customer_id": 3213,
    "customer_name": "King Inc.",
    "service_id": 12984,
    "service_name": "Writing",
    "created_at": "2015-10-16T12:39:00+02:00",
    "updated_at": "2015-10-16T12:39:00+02:00"
  }
}`

const timeEntryRequest = `{
  "time_entry": {
    "date_at": "2015-10-16",
    "minutes": 15,
    "note": "Rework description of authentication process",
    "project_id": 88309,
    "service_id": 12984
  }
}`

var timeEntryObject = domain.TimeEntry{
	Id:           domain.NewTimeEntryId(36159117),
	Minutes:      domain.NewMinutes(15),
	Date:         domain.NewLocalDate(time.Date(2015, time.October, 16, 0, 0, 0, 0, time.Local)),
	Note:         "Rework description of authentication process",
	Billable:     true,
	Locked:       false,
	Revenue:      0,
	HourlyRate:   0,
	UserId:       domain.NewUserId(211),
	UserName:     "Noah Scott",
	ProjectId:    domain.NewProjectId(88309),
	ProjectName:  "API Docs",
	CustomerId:   domain.NewCustomerId(3213),
	CustomerName: "King Inc.",
	ServiceId:    domain.NewServiceId(12984),
	ServiceName:  "Writing",
	CreatedAt:    time.Date(2015, time.October, 16, 10, 39, 0, 0, time.UTC),
	UpdatedAt:    time.Date(2015, time.October, 16, 10, 39, 0, 0, time.UTC),
}

func TestApi_TimeEntries(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(fmt.Sprintf("[%s]", timeEntryResponse)))
	}))

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	timeEntries, err := api.TimeEntries(nil)

	// then
	assert.Nil(t, err)
	assert.Equal(t, []*domain.TimeEntry{&timeEntryObject}, timeEntries)

	assert.Equal(t, http.MethodGet, rec.method)
	assert.Equal(t, "/time_entries.json", rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}

func TestApi_TimeEntries_WithQuery(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(fmt.Sprintf("[%s]", timeEntryResponse)))
	}))

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	today := domain.Today()
	query := &domain.TimeEntryQuery{
		From:      &today,
		To:        &today,
		Direction: "asc",
	}
	timeEntries, err := api.TimeEntries(query)

	// then
	assert.Nil(t, err)
	assert.Equal(t, []*domain.TimeEntry{&timeEntryObject}, timeEntries)

	assert.Equal(t, http.MethodGet, rec.method)
	assert.Equal(t, fmt.Sprintf("/time_entries.json?direction=%s&from=%s&to=%s", query.Direction, query.From, query.To), rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}

func TestApi_TimeEntry(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(timeEntryResponse))
	}))

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	timeEntry, err := api.TimeEntry(timeEntryObject.Id)

	// then
	assert.Nil(t, err)
	assert.Equal(t, &timeEntryObject, timeEntry)

	assert.Equal(t, http.MethodGet, rec.method)
	assert.Equal(t, fmt.Sprintf("/time_entries/%s.json", timeEntryObject.Id), rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}

func TestApi_CreateTimeEntry(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)

		rec.body = b
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.contentType = r.Header.Get("Content-Type")
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("Location", fmt.Sprintf("/time_entries/%s.json", timeEntryObject.Id))
		w.WriteHeader(201)
		_, _ = w.Write([]byte(timeEntryResponse))
	}))

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	command := &domain.TimeEntryCommand{
		Date:      &timeEntryObject.Date,
		Minutes:   &timeEntryObject.Minutes,
		Note:      timeEntryObject.Note,
		ProjectId: timeEntryObject.ProjectId,
		ServiceId: timeEntryObject.ServiceId,
	}
	timeEntry, err := api.CreateTimeEntry(command)

	// then
	assert.Nil(t, err)
	assert.Equal(t, &timeEntryObject, timeEntry)

	assert.Equal(t, timeEntryRequest, string(prettifyJson(rec.body, "  ")))
	assert.Equal(t, http.MethodPost, rec.method)
	assert.Equal(t, "/time_entries.json", rec.url)
	assert.Equal(t, "application/json", rec.contentType)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}

func TestApi_EditTimeEntry(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)

		rec.body = b
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.contentType = r.Header.Get("Content-Type")

		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.WriteHeader(200)
	}))

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	command := &domain.TimeEntryCommand{
		Date:      &timeEntryObject.Date,
		Minutes:   &timeEntryObject.Minutes,
		Note:      timeEntryObject.Note,
		ProjectId: timeEntryObject.ProjectId,
		ServiceId: timeEntryObject.ServiceId,
	}
	err = api.EditTimeEntry(timeEntryObject.Id, command)

	// then
	assert.Nil(t, err)

	assert.Equal(t, timeEntryRequest, string(prettifyJson(rec.body, "  ")))
	assert.Equal(t, http.MethodPatch, rec.method)
	assert.Equal(t, fmt.Sprintf("/time_entries/%s.json", timeEntryObject.Id), rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}

func TestApi_DeleteTimeEntry(t *testing.T) {
	// given
	rec := recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.url = r.RequestURI
		rec.miteKey = r.Header.Get("X-MiteApiKey")
		rec.userAgent = r.Header.Get("User-Agent")

		w.WriteHeader(200)
	}))

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	err = api.DeleteTimeEntry(timeEntryObject.Id)

	// then
	assert.Nil(t, err)

	assert.Equal(t, http.MethodDelete, rec.method)
	assert.Equal(t, fmt.Sprintf("/time_entries/%s.json", timeEntryObject.Id), rec.url)
	assert.Equal(t, testApiKey, rec.miteKey)
	assert.Equal(t, testUserAgent, rec.userAgent)
}
