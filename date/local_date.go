package date

import "time"

const ISO8601 = "2006-01-02"

type LocalDate struct {
	time time.Time
}

func Today() LocalDate {
	return From(time.Now().Local())
}

func From(t time.Time) LocalDate {
	return LocalDate{time: t}
}

func Parse(s string) (LocalDate, error) {
	t, err := time.ParseInLocation(ISO8601, s, time.Local)
	if err != nil {
		return LocalDate{}, err
	}

	return From(t), nil
}

func (d LocalDate) Add(years int, months int, days int) LocalDate {
	return LocalDate{time: d.time.AddDate(years, months, days)}
}

func (d LocalDate) String() string {
	return d.time.Format(ISO8601)
}
