package mite

import (
	"net/http"
)

const userAgent = "mite-go/0.1 (+github.com/leanovate/mite-go)"
const layout = "2006-01-02"

type MiteApi interface {
	Projects() ([]Project, error)
	Services() ([]Service, error)
	TimeEntries(params *TimeEntryParameters) ([]TimeEntry, error)
}

type defaultApi struct {
	url    string
	key    string
	client *http.Client
}

func NewMiteApi(url string, key string) MiteApi {
	return &defaultApi{url: url, key: key, client: &http.Client{}}
}
