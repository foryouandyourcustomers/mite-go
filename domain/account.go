package domain

import "strconv"

type AccountId int

func NewAccountId(i int) AccountId {
	return AccountId(i)
}

func ParseAccountId(s string) (AccountId, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return NewAccountId(i), nil
}

func (i AccountId) String() string {
	return strconv.Itoa(int(i))
}

type AccountApi interface{}
