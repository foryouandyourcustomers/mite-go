package mite

import (
	"fmt"
	"github.com/leanovate/mite-go/domain"
)

type projectResponse struct {
	Project struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Note string `json:"note"`
	} `json:"project"`
}

func (r *projectResponse) toProject() *domain.Project {
	return &domain.Project{
		Id:   fmt.Sprintf("%d", r.Project.Id),
		Name: r.Project.Name,
		Note: r.Project.Note,
	}
}

func (a *api) Projects() ([]*domain.Project, error) {
	var prs []projectResponse
	err := a.get("projects.json", &prs)
	if err != nil {
		return nil, err
	}

	var projects []*domain.Project
	for _, pr := range prs {
		projects = append(projects, pr.toProject())
	}

	return projects, nil
}
