package dtos

type UpdateTaskCommand struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

func (cmd *UpdateTaskCommand) Validate() error {

	if cmd.Title == "" && cmd.Description == "" && cmd.DueDate == "" {
		return NoParamsToUpdate
	}

	if cmd.Title != "" {
		if len(cmd.Title) > 255 {
			return TitleMaxLenExceeded
		}
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
