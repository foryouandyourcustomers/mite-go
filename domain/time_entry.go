package domain

import "time"

type TimeEntry struct {
	Id           string
	Minutes      Minutes
	Date         LocalDate
	Note         string
	Billable     bool
	Locked       bool
	Revenue      float64
	HourlyRate   int
	UserId       string
	UserName     string
	ProjectId    string
	ProjectName  string
	CustomerId   string
	CustomerName string
	ServiceId    string
	ServiceName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TimeEntryCommand struct {
	Date      *LocalDate
	Minutes   *Minutes
	Note      string
	UserId    string
	ProjectId string
	ServiceId string
	Locked    bool
}

type TimeEntryQuery struct {
	From      *LocalDate
	To        *LocalDate
	Direction string
}
