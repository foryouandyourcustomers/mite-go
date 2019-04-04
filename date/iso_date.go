package date

import "time"

const ISO8601 = "2006-01-02"

type Date struct {
	time time.Time
}

func Today() Date {
	return From(time.Now())
}

func From(time time.Time) Date {
	return Date{time: time}
}

func Parse(date string) (Date, error) {
	t, err := time.Parse(ISO8601, date)
	if err != nil {
		return Date{}, err
	}

	return From(t), nil
}

func (d Date) Add(years int, months int, days int) Date {
	return Date{time: d.time.AddDate(years, months, days)}
}

func (d Date) String() string {
	return d.time.Format(ISO8601)
}
