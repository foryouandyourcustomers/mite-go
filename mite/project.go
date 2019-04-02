package mite

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	Id   string
	Name string
	Note string
}

func (a *defaultApi) Projects() ([]Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.url, "projects.json"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-MiteApiKey", a.key)
	req.Header.Add("User-Agent", userAgent)

	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	prs := []ProjectResponse{}
	err = json.NewDecoder(res.Body).Decode(&prs)
	if err != nil {
		return nil, err
	}

	projects := []Project{}
	for _, pr := range prs {
		projects = append(projects, pr.ToProject())
	}

	return projects, nil
}

type ProjectResponse struct {
	Project struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Note string `json:"note"`
	} `json:"project"`
}

func (r ProjectResponse) ToProject() Project {
	return Project{
		Id:   fmt.Sprintf("%d", r.Project.Id),
		Name: r.Project.Name,
		Note: r.Project.Note,
	}
}
