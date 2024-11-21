package dtos

import "errors"

var (
	TitleIsRequired           = errors.New("title is required")
	TitleMaxLenExceeded       = errors.New("title cannot exceed 255 characters")
	DescriptionMaxLenExceeded = errors.New("description cannot exceed 500 characters")
	NotValidDateFormat        = errors.New("due_date must be in format YYYY-MM-DD and not less than today")
	NoParamsToUpdate          = errors.New("at least 1 parameter must be set to update")
)
