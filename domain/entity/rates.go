package entity

import "time"

type Rate struct {
	PublishedAt time.Time
	Code        string
	Value       string
}
