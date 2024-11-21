package dtos

type CreateTaskCommand struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

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
