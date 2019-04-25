package domain

import (
	"math"
	"strings"
	"time"
)

const ISO8601 = "2006-01-02"

type LocalDate struct {
	time time.Time
}

func (d *LocalDate) Before(b LocalDate) bool {
	return d.time.Before(b.time)
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

func (d LocalDate) String() string {
	return d.time.Format(ISO8601)
}

type Minutes struct {
	duration time.Duration
}

func NewMinutes(minutes int) Minutes {
	return Minutes{duration: time.Duration(minutes) * time.Minute}
}

func NewMinutesFromHours(hours int) Minutes {
	return Minutes{duration: time.Duration(hours*60) * time.Minute}
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
