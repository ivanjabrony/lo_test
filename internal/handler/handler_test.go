package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"ivanjabrony/test_lo/internal/model"
	"ivanjabrony/test_lo/internal/model/dto"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockTaskUsecase struct {
	storeFunc       func(ctx context.Context, request dto.PostTaskRequest) (int, error)
	getAllFunc      func(ctx context.Context, filter model.Filter) (dto.GetAllTasksResponse, error)
	getByTaskIdFunc func(ctx context.Context, taskId int) (dto.GetTaskByIdResponse, error)
}

func (m *MockTaskUsecase) Store(ctx context.Context, request dto.PostTaskRequest) (int, error) {
	return m.storeFunc(ctx, request)
}

func (m *MockTaskUsecase) GetAll(ctx context.Context, filter model.Filter) (dto.GetAllTasksResponse, error) {
	return m.getAllFunc(ctx, filter)
}

func (m *MockTaskUsecase) GetByTaskId(ctx context.Context, taskId int) (dto.GetTaskByIdResponse, error) {
	return m.getByTaskIdFunc(ctx, taskId)
}

type MockLogger struct {
	logs []string
}

func (m *MockLogger) Log(format string, info ...any) {
	m.logs = append(m.logs, fmt.Sprintf(format, info...))
}

func TestNewTaskHandler(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockLogger := &MockLogger{}
		mockUsecase := &MockTaskUsecase{}

		handler, err := NewTaskHandler(mockLogger, mockUsecase)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if handler == nil {
			t.Fatal("Expected handler to be initialized, got nil")
		}
	})

	t.Run("nil usecase", func(t *testing.T) {
		mockLogger := &MockLogger{}
		_, err := NewTaskHandler(mockLogger, nil)

		if err == nil {
			t.Fatal("Expected error for nil usecase, got nil")
		}

		expectedErr := "nil values in TransactionHandler constructor"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%v'", expectedErr, err)
		}
	})
}

func TestHandlePostTask(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		usecaseReturn  int
		usecaseError   error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful post",
			requestBody: dto.PostTaskRequest{
				Name:        "Test Task",
				Description: "Test Description",
				Status:      model.Created,
			},
			usecaseReturn:  42,
			usecaseError:   nil,
			expectedStatus: http.StatusOK,
			expectedBody:   42,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid",
			usecaseReturn:  -1,
			usecaseError:   nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "invalid data in task"},
		},
		{
			name: "invalid task status",
			requestBody: dto.PostTaskRequest{
				Name:        "Test Task",
				Description: "Test Description",
				Status:      "invalid-status",
			},
			usecaseReturn:  -1,
			usecaseError:   nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "invalid data in task"},
		},
		{
			name: "usecase error",
			requestBody: dto.PostTaskRequest{
				Name:        "Test Task",
				Description: "Test Description",
				Status:      model.Done,
			},
			usecaseReturn:  -1,
			usecaseError:   errors.New("usecase error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "failed to store task"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			mockUsecase := &MockTaskUsecase{
				storeFunc: func(ctx context.Context, request dto.PostTaskRequest) (int, error) {
					return tt.usecaseReturn, tt.usecaseError
				},
			}

			handler, _ := NewTaskHandler(mockLogger, mockUsecase)

			var reqBody []byte
			var err error
			reqBody, err = json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandlePostTask(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectedStatus == http.StatusOK {
				var id int
				if err := json.NewDecoder(resp.Body).Decode(&id); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}
				if id != tt.usecaseReturn {
					t.Errorf("Expected ID %d, got %d", tt.usecaseReturn, id)
				}
			} else {
				var errorResponse map[string]string
				if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if errorResponse["error"] != tt.expectedBody.(map[string]string)["error"] {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedBody.(map[string]string)["error"], errorResponse["error"])
				}
			}
		})
	}
}

func TestHandleGetAllTasks(t *testing.T) {
	now := time.Now()
	testTasks := []model.Task{
		{Id: 1, Name: "Task 1", Description: "Desc 1", Status: model.Created, CreatedAt: now},
		{Id: 2, Name: "Task 2", Description: "Desc 2", Status: model.InProgress, CreatedAt: now},
	}

	tests := []struct {
		name           string
		queryParams    map[string]string
		usecaseReturn  dto.GetAllTasksResponse
		usecaseError   error
		expectedStatus int
		expectedLength int
	}{
		{
			name:           "successful get all",
			queryParams:    map[string]string{},
			usecaseReturn:  dto.GetAllTasksResponse{Amount: 2, Tasks: testTasks},
			usecaseError:   nil,
			expectedStatus: http.StatusOK,
			expectedLength: 2,
		},
		{
			name:           "filtered get all",
			queryParams:    map[string]string{"status": "created"},
			usecaseReturn:  dto.GetAllTasksResponse{Amount: 1, Tasks: []model.Task{testTasks[0]}},
			usecaseError:   nil,
			expectedStatus: http.StatusOK,
			expectedLength: 1,
		},
		{
			name:           "invalid filter",
			queryParams:    map[string]string{"status": "invalid"},
			usecaseReturn:  dto.GetAllTasksResponse{},
			usecaseError:   nil,
			expectedStatus: http.StatusBadRequest,
			expectedLength: 0,
		},
		{
			name:           "usecase error",
			queryParams:    map[string]string{},
			usecaseReturn:  dto.GetAllTasksResponse{},
			usecaseError:   errors.New("usecase error"),
			expectedStatus: http.StatusInternalServerError,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			mockUsecase := &MockTaskUsecase{
				getAllFunc: func(ctx context.Context, filter model.Filter) (dto.GetAllTasksResponse, error) {
					return tt.usecaseReturn, tt.usecaseError
				},
			}

			handler, _ := NewTaskHandler(mockLogger, mockUsecase)

			req := httptest.NewRequest("GET", "/tasks", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()

			handler.HandleGetAllTasks(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectedStatus == http.StatusOK {
				var response dto.GetAllTasksResponse
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if response.Amount != tt.expectedLength {
					t.Errorf("Expected %d tasks, got %d", tt.expectedLength, response.Amount)
				}
			} else if tt.expectedStatus == http.StatusBadRequest {
				var errorResponse map[string]string
				if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}

				if errorResponse["error"] != "invalid task status in filter" {
					t.Errorf("Unexpected error message: %s", errorResponse["error"])
				}
			}
		})
	}
}

func TestHandleGetTaskById(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		path           string
		usecaseReturn  dto.GetTaskByIdResponse
		usecaseError   error
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful get by id",
			path: "/tasks/42",
			usecaseReturn: dto.GetTaskByIdResponse{
				Id:          42,
				Name:        "Test Task",
				Description: "Test Description",
				Status:      model.Done,
				CreatedAt:   now,
			},
			usecaseError:   nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing task id",
			path:           "/tasks/",
			usecaseReturn:  dto.GetTaskByIdResponse{},
			usecaseError:   nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "task_id wasn't provided",
		},
		{
			name:           "invalid task id",
			path:           "/tasks/invalid",
			usecaseReturn:  dto.GetTaskByIdResponse{},
			usecaseError:   nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid task_id parameter",
		},
		{
			name:           "task not found",
			path:           "/tasks/99",
			usecaseReturn:  dto.GetTaskByIdResponse{},
			usecaseError:   errors.New("not found"),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to retrieve task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			mockUsecase := &MockTaskUsecase{
				getByTaskIdFunc: func(ctx context.Context, taskId int) (dto.GetTaskByIdResponse, error) {
					if tt.path == "/tasks/42" && taskId != 42 {
						t.Errorf("Task ID mismatch. Expected 42, got %d", taskId)
					}
					return tt.usecaseReturn, tt.usecaseError
				},
			}

			handler, _ := NewTaskHandler(mockLogger, mockUsecase)

			req := httptest.NewRequest("GET", tt.path, nil)

			if strings.HasPrefix(tt.path, "/tasks/") {
				parts := strings.Split(tt.path, "/")
				if len(parts) > 2 {
					req.SetPathValue("task_id", parts[2])
				}
			}

			w := httptest.NewRecorder()

			handler.HandleGetTaskById(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if tt.expectedStatus == http.StatusOK {
				var response dto.GetTaskByIdResponse
				if err := json.Unmarshal(bodyBytes, &response); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if response.Id != tt.usecaseReturn.Id {
					t.Errorf("Expected task ID %d, got %d", tt.usecaseReturn.Id, response.Id)
				}
			} else {
				var errorResponse map[string]string
				if err := json.Unmarshal(bodyBytes, &errorResponse); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}

				if errorResponse["error"] != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, errorResponse["error"])
				}
			}
		})
	}
}
