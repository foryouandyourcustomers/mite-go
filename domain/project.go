package domain

type Project struct {
	Id   string
	Name string
	Note string
}

type ProjectApi interface {
	Projects() ([]*Project, error)
}
