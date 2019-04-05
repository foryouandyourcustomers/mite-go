package date

import "time"

const ISO8601 = "2006-01-02"

type LocalDate struct {
	time time.Time
}

func NewLocalDate(t time.Time) LocalDate {
	return LocalDate{time: t}
}

func Today() LocalDate {
	return NewLocalDate(time.Now().Local())
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
