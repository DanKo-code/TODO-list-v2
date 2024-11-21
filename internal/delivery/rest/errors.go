package rest

import "errors"

var (
	TaskIdIsRequired                 = errors.New("task id is required")
	InvalidIdFormat                  = errors.New("id must be on uuid format")
	NoParamsToUpdate                 = errors.New("at least 1 parameter must be set to update")
	NoParamsToCreate                 = errors.New("at least 1 parameter(title) must be set to create")
	NoParamsToChangeCompletionStatus = errors.New("completion status is required")
)
