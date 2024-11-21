package dtos

import "time"

func isValidDate(date string) bool {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}

	today := time.Now().Truncate(24 * time.Hour)

	return !parsedDate.Before(today)
}
