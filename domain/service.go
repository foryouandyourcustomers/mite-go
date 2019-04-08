package domain

import "strconv"

type ServiceId int

func NewServiceId(i int) ServiceId {
	return ServiceId(i)
}

func ParseServiceId(s string) (ServiceId, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return NewServiceId(i), nil
}

func (i ServiceId) String() string {
	return strconv.Itoa(int(i))
}

type Service struct {
	Id   ServiceId
	Name string
	Note string
}

type ServiceApi interface {
	Services() ([]*Service, error)
}
