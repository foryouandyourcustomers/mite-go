package domain

import (
	"strconv"
	"time"
)

const (
	SORT_BY_USER     = "user"
	SORT_BY_CUSTOMER = "customer"
	SORT_BY_PROJECT  = "project"
	SORT_BY_SERVICE  = "service"
	SORT_BY_NOTE     = "note"
	SORT_BY_TIME     = "minutes"
	SORT_BY_REVENUE  = "revenue"

	SORT_DIRECTION_ASC  = "asc"
	SORT_DIRECTION_DESC = "dsc"
)

type TimeEntryId int

func NewTimeEntryId(i int) TimeEntryId {
	return TimeEntryId(i)
}

func ParseTimeEntryId(s string) (TimeEntryId, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return NewTimeEntryId(i), nil
}

func (i TimeEntryId) String() string {
	return strconv.Itoa(int(i))
}

type TimeEntry struct {
	Id           TimeEntryId
	Minutes      Minutes
	Date         LocalDate
	Note         string
	Billable     bool
	Locked       bool
	Revenue      float64
	HourlyRate   int
	UserId       UserId
	UserName     string
	ProjectId    ProjectId
	ProjectName  string
	CustomerId   CustomerId
	CustomerName string
	ServiceId    ServiceId
	ServiceName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TimeEntryCommand struct {
	Date      *LocalDate
	Minutes   *Minutes
	Note      string
	UserId    UserId
	ProjectId ProjectId
	ServiceId ServiceId
	Locked    bool
}

type TimeEntryQuery struct {
	At         string
	From       *LocalDate
	To         *LocalDate
	Direction  string
	ServiceId  ServiceId
	UserId     UserId
	CustomerId CustomerId
	Sort       string
}

type TimeEntryApi interface {
	TimeEntries(query *TimeEntryQuery) ([]*TimeEntry, error)
	TimeEntry(id TimeEntryId) (*TimeEntry, error)
	CreateTimeEntry(command *TimeEntryCommand) (*TimeEntry, error)
	EditTimeEntry(id TimeEntryId, command *TimeEntryCommand) error
	DeleteTimeEntry(id TimeEntryId) error
}
