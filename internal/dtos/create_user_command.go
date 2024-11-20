package dtos

import (
	"errors"
	"time"
)

type CreateTaskCommand struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

var (
	TitleIsRequired           = errors.New("title is required")
	TitleMaxLenExceeded       = errors.New("title cannot exceed 255 characters")
	DescriptionMaxLenExceeded = errors.New("description cannot exceed 500 characters")
	NotValidDateFormat        = errors.New("due_date must be in format YYYY-MM-DD and bigger then today")
)

func (cmd *CreateTaskCommand) Validate() error {
	if cmd.Title == "" {
		return TitleIsRequired
	}
	if len(cmd.Title) > 255 {
		return TitleMaxLenExceeded
	}

	if cmd.Description != "" {
		if len(cmd.Description) > 500 {
			return DescriptionMaxLenExceeded
		}
	}

	if cmd.DueDate != "" {
		if !isValidDate(cmd.DueDate) {
			return NotValidDateFormat
		}
	}

	return nil
}

func isValidDate(date string) bool {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}

	today := time.Now().Truncate(24 * time.Hour)

	return parsedDate.After(today)
}
