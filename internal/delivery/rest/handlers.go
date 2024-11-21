package rest

import (
	"github.com/DanKo-code/TODO-list/internal/dtos"
	"github.com/DanKo-code/TODO-list/internal/usecase"
	"net/http"
)

var (
	NoBody = "EOF"
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

		if err.Error() == NoBody {
			WriteErrToResponseBody(w, NoParamsToCreate, http.StatusNotFound)
			return
		}

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
		WriteErrToResponseBody(w, err, http.StatusInternalServerError)
		return
	}

	WriteToResponseBody(w, task)
}

func (h *Handlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := h.useCase.GetTasks(ctx)
	if err != nil {
		WriteErrToResponseBody(w, err, http.StatusInternalServerError)
		return
	}

	WriteToResponseBody(w, tasks)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskId, ok := ctx.Value("id").(string)
	if !ok || !isValidUUID(taskId) {
		WriteErrToResponseBody(w, InvalidIdFormat, http.StatusBadRequest)
		return
	}

	cmd := dtos.UpdateTaskCommand{}
	err := ReadFromRequestBody(r, &cmd)
	if err != nil {

		if err.Error() == NoBody {
			WriteErrToResponseBody(w, NoParamsToUpdate, http.StatusNotFound)
			return
		}

		WriteErrToResponseBody(w, err, http.StatusBadRequest)
		return
	}

	err = cmd.Validate()
	if err != nil {
		WriteErrToResponseBody(w, err, http.StatusBadRequest)
		return
	}

	utask, err := h.useCase.UpdateTask(ctx, taskId, &cmd)
	if err != nil {
		WriteErrToResponseBody(w, err, http.StatusInternalServerError)
		return
	}

	WriteToResponseBody(w, utask)
}
