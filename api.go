package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
