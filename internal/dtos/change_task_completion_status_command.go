package dtos

type ChangeTaskCompletionStatusCommand struct {
	Completed *bool `json:"completed"`
}

func (cmd *ChangeTaskCompletionStatusCommand) Validate() error {
	if cmd.Completed == nil {
		return CompletedIsRequired
	}

	return nil
}
