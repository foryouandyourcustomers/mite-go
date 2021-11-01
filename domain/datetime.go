package domain

import (
	"math"
	"strings"
	"time"
)

const (
	ISO8601       = "2006-01-02"
	AT_TODAY      = "today"
	AT_YESTERDAY  = "yesterday"
	AT_THIS_WEEK  = "this_week"
	AT_THIS_MONTH = "this_month"
	AT_THIS_YEAR  = "this_year"
	AT_LAST_WEEK  = "last_week"
	AT_LAST_MONTH = "last_month"
	AT_LAST_YEAR  = "last_year"
)

type LocalDate struct {
	time time.Time
}

func NewLocalDate(t time.Time) LocalDate {
	return LocalDate{time: t}
}

func Today() LocalDate {
	return NewLocalDate(time.Now().Local())
}

func ThisYear() int {
	return time.Now().Year()
}

func ParseLocalDate(s string) (LocalDate, error) {
	t, err := time.ParseInLocation(ISO8601, s, time.Local)
	if err != nil {
		return LocalDate{}, err
	}

	return NewLocalDate(t), nil
}

func (d LocalDate) Add(years int, months int, days int) LocalDate {
	return LocalDate{time: d.time.AddDate(years, months, days)}
}

func (d LocalDate) Before(b LocalDate) bool {
	return d.time.Before(b.time)
}

func (d LocalDate) Unix() int64 {
	return d.time.Unix()
}

func (d LocalDate) String() string {
	return d.time.Format(ISO8601)
}

func (d LocalDate) Month() string {
	return d.time.Month().String()
}

func (d LocalDate) Day() int {
	return d.time.Day()
}

func (d LocalDate) Year() int {
	return d.time.Year()
}

type Minutes struct {
	duration time.Duration
}

func NewMinutes(minutes int) Minutes {
	return Minutes{duration: time.Duration(minutes) * time.Minute}
}

func NewMinutesFromHours(hours int) Minutes {
	return Minutes{duration: time.Duration(hours) * time.Hour}
}

func ParseMinutes(s string) (Minutes, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return Minutes{}, err
	}

	return Minutes{duration: d.Round(time.Minute)}, nil
}

func (m Minutes) Value() int {
	return int(math.Min(m.duration.Minutes(), math.MaxInt32))
}

func (m Minutes) String() string {
	return strings.TrimSuffix(m.duration.String(), "0s")
}

func MinutesAsDays(minutes int, workingDayInHours float64) float64 {
	return float64(minutes) / 60 / workingDayInHours
}
