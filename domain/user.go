package domain

import "strconv"

const (
	currentUserId     = -1
	currentUserString = "current"
)
const CurrentUser = UserId(currentUserId)

type UserId int

func NewUserId(i int) UserId {
	return UserId(i)
}

func ParseUserId(s string) (UserId, error) {
	if s == currentUserString {
		return CurrentUser, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return NewUserId(i), nil
}

func (i UserId) String() string {
	if i == currentUserId {
		return currentUserString
	}
	return strconv.Itoa(int(i))
}

type UserApi interface{}
