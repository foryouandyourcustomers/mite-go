package domain

type Service struct {
	Id   string
	Name string
	Note string
}

type ServiceApi interface {
	Services() ([]*Service, error)
}
