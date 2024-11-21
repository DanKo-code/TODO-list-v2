package dtos

type ChangeTaskCompletionStatusCommand struct {
	Completed bool `json:"completed"`
}
