package mite

import (
	"fmt"
)

type Project struct {
	Id   string
	Name string
	Note string
}

type projectResponse struct {
	Project struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Note string `json:"note"`
	} `json:"project"`
}

func (r *projectResponse) ToProject() *Project {
	return &Project{
		Id:   fmt.Sprintf("%d", r.Project.Id),
		Name: r.Project.Name,
		Note: r.Project.Note,
	}
}

func (a *miteApi) Projects() ([]*Project, error) {
	prs := []projectResponse{}
	err := a.get("projects.json", &prs)
	if err != nil {
		return nil, err
	}

	projects := []*Project{}
	for _, pr := range prs {
		projects = append(projects, pr.ToProject())
	}

	return projects, nil
}
