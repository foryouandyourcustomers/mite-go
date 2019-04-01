package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Project struct {
	ProjectBody GenericBody `json:"project"`
}

type GenericBody struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Note string `json:"note"`
}

type Service struct {
	ServiceBody GenericBody `json:"service"`
}

type TimeEntry struct {
	TimeEntryBody TimeEntryBody `json:"time_entry"`
}

type TimeEntryBody struct {
	GenericBody
	Minutes     int    `json:"minutes"`
	Date        string `json:"date_at"`
	ProjectName string `json:"project_name"`
	ServiceName string `json:"service_name"`
}

func apiGetProjects() (projects []Project) {
	client := &http.Client{}
	req := buildGetRequest("/projects.json")
	res, err := client.Do(req)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&projects)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	return projects
}

func buildGetRequest(path string) *http.Request {
	url := fmt.Sprintf("%s%s", configGetApiUrl(), path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	req.Header.Add("X-MiteApiKey", configGetApiKey())
	req.Header.Add("User-Agent", "mite-go/0.1 (+github.com/phiros/mite-go)")

	return req
}

func apiGetServices() (services []Service) {
	client := &http.Client{}
	req := buildGetRequest("/services.json")
	res, err := client.Do(req)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&services)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	return services
}

func apiGetEntries() (timeEntries []TimeEntry) {
	client := &http.Client{}

	toDate := time.Now().Format("2006-01-02")
	fromDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	requestPathAndParams := fmt.Sprintf("/time_entries.json?from=%s&to=%s&direction=asc", fromDate, toDate)
	req := buildGetRequest(requestPathAndParams)
	res, err := client.Do(req)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&timeEntries)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	return timeEntries
}
