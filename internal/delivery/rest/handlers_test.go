package rest

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DanKo-code/TODO-list/internal/dtos"
	internalErrors "github.com/DanKo-code/TODO-list/internal/errors"
	"github.com/DanKo-code/TODO-list/internal/models"
	"github.com/DanKo-code/TODO-list/internal/usecase/task_usecase"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTaskHandler(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		mockCreateTaskFunc func(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "success",
			requestBody: `{"title":"Test Task","description":"This is a test task"," due_date":"2024-11-22"}`,
			mockCreateTaskFunc: func(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error) {
				return &models.Task{Id: "a495465c-d177-48e1-8954-516bba76d541", Title: "Test Task", Description: "This is a test task", DueDate: "2024-11-22", Overdue: false, Completed: false}, nil
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id":"a495465c-d177-48e1-8954-516bba76d541","title":"Test Task","description":"This is a test task","due_date":"2024-11-22","overdue":false,"completed":false}`,
		},
		{
			name:        "no body",
			requestBody: ``,
			mockCreateTaskFunc: func(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error) {
				return nil, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%s"}`, NoParamsToCreate.Error()),
		},
		{
			name:        "cannot unmarshal",
			requestBody: `{"title":2,"description":"This is a test task"," due_date":"2024-11-22"}`,
			mockCreateTaskFunc: func(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error) {
				return nil, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"json: cannot unmarshal number into Go struct field CreateTaskCommand.title of type string"}`),
		},
		{
			name:        "validation error",
			requestBody: `{"description":"This is a test task"," due_date":"2024-11-22"}`,
			mockCreateTaskFunc: func(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error) {
				return nil, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%s"}`, dtos.TitleIsRequired.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &task_usecase.MockTaskUseCase{
				CreateTaskFunc: tt.mockCreateTaskFunc,
			}

			h := NewHandlers(mockUseCase)

			req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.CreateTask(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, resp.StatusCode)
			}

			if tt.expectedResponse != "" {
				var buf bytes.Buffer
				buf.ReadFrom(resp.Body)

				if strings.TrimSpace(buf.String()) != tt.expectedResponse {
					t.Errorf("expected %s, got %s", tt.expectedResponse, buf.String())
				}
			}

		})
	}
}

func TestGetTasksHandler(t *testing.T) {
	tests := []struct {
		name               string
		mockGetTasksFunc   func(ctx context.Context) ([]*models.Task, error)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "success",
			mockGetTasksFunc: func(ctx context.Context) ([]*models.Task, error) {
				return []*models.Task{
					{Id: "a495465c-d177-48e1-8954-516bba76d541", Title: "Test Task", Description: "This is a test task", DueDate: "2024-11-22", Overdue: false, Completed: false},
				}, nil
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`[{"id":"a495465c-d177-48e1-8954-516bba76d541","title":"Test Task","description":"This is a test task","due_date":"2024-11-22","overdue":false,"completed":false}]`),
		},
		{
			name: "success",
			mockGetTasksFunc: func(ctx context.Context) ([]*models.Task, error) {
				return nil, nil
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`[]`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &task_usecase.MockTaskUseCase{
				GetTaskFunc: tt.mockGetTasksFunc,
			}
			h := NewHandlers(mockUseCase)
			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			w := httptest.NewRecorder()
			h.GetTasks(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, resp.StatusCode)
			}
			if tt.expectedResponse != "" {
				var buf bytes.Buffer
				buf.ReadFrom(resp.Body)

				if strings.TrimSpace(buf.String()) != tt.expectedResponse {
					t.Errorf("expected %s, got %s", tt.expectedResponse, buf.String())
				}
			}
		})
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		requestBody        string
		mockUpdateTaskFunc func(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) (*models.Task, error)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "success",
			id:          "a495465c-d177-48e1-8954-516bba76d541",
			requestBody: `{"title":"Test Task","description":"This is a test task","due_date":"2024-11-28"}`,
			mockUpdateTaskFunc: func(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) (*models.Task, error) {
				return &models.Task{
					Id:          "a495465c-d177-48e1-8954-516bba76d541",
					Title:       "Test Task",
					Description: "This is a test task",
					DueDate:     "2024-11-22",
					Overdue:     false,
					Completed:   false,
				}, nil
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"id":"a495465c-d177-48e1-8954-516bba76d541","title":"Test Task","description":"This is a test task","due_date":"2024-11-22","overdue":false,"completed":false}`),
		},
		{
			name:        "invalid id format",
			id:          "1",
			requestBody: `{"title":"Test Task","description":"This is a test task","due_date":"2024-11-28"}`,
			mockUpdateTaskFunc: func(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) (*models.Task, error) {
				return &models.Task{
					Id:          "a495465c-d177-48e1-8954-516bba76d541",
					Title:       "Test Task",
					Description: "This is a test task",
					DueDate:     "2024-11-22",
					Overdue:     false,
					Completed:   false,
				}, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   fmt.Sprintf(`{"error":"%s"}`, InvalidIdFormat),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &task_usecase.MockTaskUseCase{
				UpdateTaskFunc: tt.mockUpdateTaskFunc,
			}
			h := NewHandlers(mockUseCase)

			ctx := context.Background()

			ctx = context.WithValue(ctx, "id", tt.id)

			req := httptest.NewRequestWithContext(ctx, http.MethodPut, "/tasks/a495465c-d177-48e1-8954-516bba76d541", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			h.UpdateTask(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, resp.StatusCode)
			}
			if tt.expectedResponse != "" {
				var buf bytes.Buffer
				buf.ReadFrom(resp.Body)

				if strings.TrimSpace(buf.String()) != tt.expectedResponse {
					t.Errorf("expected %s, got %s", tt.expectedResponse, buf.String())
				}
			}
		})
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		mockDeleteTaskFunc func(ctx context.Context, id string) error
		expectedStatusCode int
	}{
		{
			name: "success",
			id:   "a495465c-d177-48e1-8954-516bba76d541",
			mockDeleteTaskFunc: func(ctx context.Context, id string) error {
				return nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "task not found",
			id:   "a495465c-d177-48e1-8954-516bba76d541",
			mockDeleteTaskFunc: func(ctx context.Context, id string) error {
				return internalErrors.TaskNotFound
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &task_usecase.MockTaskUseCase{
				DeleteTaskFunc: tt.mockDeleteTaskFunc,
			}
			h := NewHandlers(mockUseCase)

			ctx := context.Background()

			ctx = context.WithValue(ctx, "id", tt.id)

			req := httptest.NewRequestWithContext(ctx, http.MethodDelete, "/tasks/a495465c-d177-48e1-8954-516bba76d541", nil)
			w := httptest.NewRecorder()
			h.DeleteTask(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}

func TestChangeTaskCompletionHandler(t *testing.T) {
	tests := []struct {
		name                         string
		id                           string
		requestBody                  string
		mockChangeTaskCompletionFunc func(ctx context.Context, id string, completionStatus bool) (*models.Task, error)
		expectedStatusCode           int
		expectedResponse             string
	}{
		{
			name:        "success",
			id:          "a495465c-d177-48e1-8954-516bba76d541",
			requestBody: `{"completed":true}`,
			mockChangeTaskCompletionFunc: func(ctx context.Context, id string, completionStatus bool) (*models.Task, error) {
				return &models.Task{
					Id:          "a495465c-d177-48e1-8954-516bba76d541",
					Title:       "Test Task",
					Description: "This is a test task",
					DueDate:     "2024-11-22",
					Overdue:     false,
					Completed:   true,
				}, nil
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"id":"a495465c-d177-48e1-8954-516bba76d541","title":"Test Task","description":"This is a test task","due_date":"2024-11-22","overdue":false,"completed":true}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &task_usecase.MockTaskUseCase{
				ChangeTaskCompletionStatusFunc: tt.mockChangeTaskCompletionFunc,
			}
			h := NewHandlers(mockUseCase)

			ctx := context.Background()

			ctx = context.WithValue(ctx, "id", tt.id)

			req := httptest.NewRequestWithContext(ctx, http.MethodPatch, "/tasks/a495465c-d177-48e1-8954-516bba76d541/complete", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			h.ChangeTaskCompletionStatus(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, resp.StatusCode)
			}
			if tt.expectedResponse != "" {
				var buf bytes.Buffer
				buf.ReadFrom(resp.Body)

				if strings.TrimSpace(buf.String()) != tt.expectedResponse {
					t.Errorf("expected %s, got %s", tt.expectedResponse, buf.String())
				}
			}
		})
	}
}
