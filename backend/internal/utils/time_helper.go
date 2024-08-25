package utils

import (
	"time"
)

func ParseDate(dateStr *string) (*time.Time, error) {
	if dateStr == nil {
		return nil, nil
	}
	parsedTime, err := time.Parse(time.RFC3339, *dateStr)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}

func FormatDate(date *time.Time) *string {
	if date == nil {
		return nil
	}
	formatted := date.Format(time.RFC3339)
	return &formatted
}
