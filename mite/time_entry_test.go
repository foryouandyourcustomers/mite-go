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
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(fmt.Sprintf("[%s]", timeEntryResponse)).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	timeEntries, err := api.TimeEntries(nil)

	// then
	assert.Nil(t, err)
	assert.Equal(t, []*domain.TimeEntry{&timeEntryObject}, timeEntries)

	assert.Equal(t, http.MethodGet, rec.RequestMethod())
	assert.Equal(t, "/time_entries.json", rec.RequestURI())
	assert.Empty(t, rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Empty(t, rec.RequestBody())
}

func TestApi_TimeEntries_WithQuery(t *testing.T) {
	// given
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(fmt.Sprintf("[%s]", timeEntryResponse)).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

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

	assert.Equal(t, http.MethodGet, rec.RequestMethod())
	assert.Equal(t, fmt.Sprintf("/time_entries.json?direction=%s&from=%s&to=%s", query.Direction, query.From, query.To), rec.RequestURI())
	assert.Empty(t, rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Empty(t, rec.RequestBody())
}

func TestApi_TimeEntry(t *testing.T) {
	// given
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseBody(timeEntryResponse).
		ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	timeEntry, err := api.TimeEntry(timeEntryObject.Id)

	// then
	assert.Nil(t, err)
	assert.Equal(t, &timeEntryObject, timeEntry)

	assert.Equal(t, http.MethodGet, rec.RequestMethod())
	assert.Equal(t, fmt.Sprintf("/time_entries/%s.json", timeEntryObject.Id), rec.RequestURI())
	assert.Empty(t, rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Empty(t, rec.RequestBody())
}

func TestApi_CreateTimeEntry(t *testing.T) {
	// given
	rec := NewRecorder().
		ResponseContentType("application/json; charset=utf-8").
		ResponseLocation(fmt.Sprintf("/time_entries/%s.json", timeEntryObject.Id)).
		ResponseBody(timeEntryResponse).
		ResponseStatus(201)
	srv := httptest.NewServer(rec.Handler())

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

	assert.Equal(t, http.MethodPost, rec.RequestMethod())
	assert.Equal(t, "/time_entries.json", rec.RequestURI())
	assert.Equal(t, "application/json", rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Equal(t, timeEntryRequest, string(prettifyJson(rec.RequestBody(), "  ")))
}

func TestApi_EditTimeEntry(t *testing.T) {
	// given
	rec := NewRecorder().ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

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

	assert.Equal(t, http.MethodPatch, rec.RequestMethod())
	assert.Equal(t, fmt.Sprintf("/time_entries/%s.json", timeEntryObject.Id), rec.RequestURI())
	assert.Equal(t, "application/json", rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Equal(t, timeEntryRequest, string(prettifyJson(rec.RequestBody(), "  ")))
}

func TestApi_DeleteTimeEntry(t *testing.T) {
	// given
	rec := NewRecorder().ResponseStatus(200)
	srv := httptest.NewServer(rec.Handler())

	defer srv.Close()

	api, err := mite.NewApi(srv.URL, testApiKey, testClientVersion)
	assert.Nil(t, err)

	// when
	err = api.DeleteTimeEntry(timeEntryObject.Id)

	// then
	assert.Nil(t, err)

	assert.Equal(t, http.MethodDelete, rec.RequestMethod())
	assert.Equal(t, fmt.Sprintf("/time_entries/%s.json", timeEntryObject.Id), rec.RequestURI())
	assert.Empty(t, rec.RequestContentType())
	assert.Equal(t, testUserAgent, rec.RequestUserAgent())
	assert.Equal(t, testApiKey, rec.RequestMiteKey())
	assert.Empty(t, rec.RequestBody())
}
