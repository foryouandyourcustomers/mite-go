package domain

import "strconv"

type UserId int

func NewUserId(i int) UserId {
	return UserId(i)
}

func ParseUserId(s string) (UserId, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return NewUserId(i), nil
}

func (i UserId) String() string {
	return strconv.Itoa(int(i))
}

type UserApi interface{}
