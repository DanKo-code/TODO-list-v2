package rest

import (
	"github.com/DanKo-code/TODO-list/internal/dtos"
	"github.com/DanKo-code/TODO-list/internal/usecase"
	"net/http"
)

type Handlers struct {
	useCase usecase.TaskUseCase
}

func NewHandlers(useCase usecase.TaskUseCase) *Handlers {
	return &Handlers{useCase}
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cmd := dtos.CreateTaskCommand{}
	err := ReadFromRequestBody(r, &cmd)
	if err != nil {
		WriteErrToResponseBody(w, err, http.StatusBadRequest)
		return
	}

	err = cmd.Validate()
	if err != nil {
		WriteErrToResponseBody(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.useCase.CreateTask(ctx, &cmd)
	if err != nil {

	}

	WriteToResponseBody(w, task)
}
