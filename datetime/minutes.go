package datetime

import (
	"math"
	"time"
)

type Minutes struct {
	duration time.Duration
}

func NewMinutes(i int) Minutes {
	return Minutes{duration: time.Duration(i) * time.Minute}
}

func ParseMinutes(s string) (Minutes, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return Minutes{}, err
	}

	return Minutes{duration: d.Truncate(time.Minute)}, nil
}

func (m Minutes) Value() int {
	return int(math.Min(m.duration.Minutes(), math.MaxInt32))
}

func (m Minutes) String() string {
	return m.duration.String()
}
