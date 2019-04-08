package domain

import "strconv"

type ProjectId int

func NewProjectId(i int) ProjectId {
	return ProjectId(i)
}

func ParseProjectId(s string) (ProjectId, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return NewProjectId(i), nil
}

func (i ProjectId) String() string {
	return strconv.Itoa(int(i))
}

type Project struct {
	Id   ProjectId
	Name string
	Note string
}

type ProjectApi interface {
	Projects() ([]*Project, error)
}
