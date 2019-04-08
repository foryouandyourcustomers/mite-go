package domain

import "strconv"

type CustomerId int

func NewCustomerId(i int) CustomerId {
	return CustomerId(i)
}

func ParseCustomerId(s string) (CustomerId, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return NewCustomerId(i), nil
}

func (i CustomerId) String() string {
	return strconv.Itoa(int(i))
}

type CustomerApi interface{}
