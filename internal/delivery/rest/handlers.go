package rest

import (
	"errors"
	"github.com/DanKo-code/TODO-list/internal/dtos"
	internalErrors "github.com/DanKo-code/TODO-list/internal/errors"
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
			WriteErrToResponseBody(w, NoParamsToCreate, http.StatusBadRequest)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteToResponseBody(w, task)
}

func (h *Handlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := h.useCase.GetTasks(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		WriteToResponseBody(w, []struct{}{})
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

		if errors.Is(err, internalErrors.TaskNotFound) {
			WriteErrToResponseBody(w, err, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteToResponseBody(w, utask)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskId, ok := ctx.Value("id").(string)
	if !ok || !isValidUUID(taskId) {
		WriteErrToResponseBody(w, InvalidIdFormat, http.StatusBadRequest)
		return
	}

	err := h.useCase.DeleteTask(ctx, taskId)
	if err != nil {
		if errors.Is(err, internalErrors.TaskNotFound) {
			WriteErrToResponseBody(w, err, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ChangeTaskCompletionStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskId, ok := ctx.Value("id").(string)
	if !ok || !isValidUUID(taskId) {
		WriteErrToResponseBody(w, InvalidIdFormat, http.StatusBadRequest)
		return
	}

	cmd := dtos.ChangeTaskCompletionStatusCommand{}
	err := ReadFromRequestBody(r, &cmd)
	if err != nil {

		if err.Error() == NoBody {
			WriteErrToResponseBody(w, NoParamsToChangeCompletionStatus, http.StatusNotFound)
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

	updatedTask, err := h.useCase.ChangeTaskCompletionStatus(ctx, taskId, *cmd.Completed)
	if err != nil {

		if errors.Is(err, internalErrors.TaskNotFound) {
			WriteErrToResponseBody(w, err, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteToResponseBody(w, updatedTask)
}
