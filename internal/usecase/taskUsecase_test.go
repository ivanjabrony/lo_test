package usecase

import (
	"context"
	"errors"
	"fmt"
	"ivanjabrony/test_lo/internal/model"
	"ivanjabrony/test_lo/internal/model/dto"
	"testing"
	"time"
)

// MockTaskStorage is a mock implementation of TaskStorage for testing
type MockTaskStorage struct {
	storeFunc       func(ctx context.Context, task model.Task) (int, error)
	getAllFunc      func(ctx context.Context, filter model.Filter) ([]model.Task, error)
	getByTaskIdFunc func(ctx context.Context, taskId int) (*model.Task, error)
}

func (m *MockTaskStorage) Store(ctx context.Context, task model.Task) (int, error) {
	return m.storeFunc(ctx, task)
}

func (m *MockTaskStorage) GetAll(ctx context.Context, filter model.Filter) ([]model.Task, error) {
	return m.getAllFunc(ctx, filter)
}

func (m *MockTaskStorage) GetByTaskId(ctx context.Context, taskId int) (*model.Task, error) {
	return m.getByTaskIdFunc(ctx, taskId)
}

type MockLogger struct {
	logs []string
}

func (m *MockLogger) Log(format string, info ...any) {
	m.logs = append(m.logs, fmt.Sprintf(format, info...))
}

func TestNewTaskUsecase(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockLogger := &MockLogger{}
		mockStorage := &MockTaskStorage{}

		usecase, err := NewTaskUsecase(mockLogger, mockStorage)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if usecase == nil {
			t.Fatal("Expected usecase to be initialized, got nil")
		}

		if len(mockLogger.logs) != 1 {
			t.Fatalf("Expected 1 log message, got %d", len(mockLogger.logs))
		}

		expectedLog := "Created TaskUsecase successfully"
		if mockLogger.logs[0] != expectedLog {
			t.Errorf("Expected log message '%s', got '%s'", expectedLog, mockLogger.logs[0])
		}
	})

	t.Run("nil storage", func(t *testing.T) {
		mockLogger := &MockLogger{}
		_, err := NewTaskUsecase(mockLogger, nil)

		if err == nil {
			t.Fatal("Expected error for nil storage, got nil")
		}

		expectedErr := "nil values in TaskUsecase constructor"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%v'", expectedErr, err)
		}
	})
}

func TestStore(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		request       dto.PostTaskRequest
		storageReturn int
		storageError  error
		wantId        int
		wantError     error
	}{
		{
			name: "successful storage",
			request: dto.PostTaskRequest{
				Name:        "Test Task",
				Description: "Test Description",
				Status:      model.Created,
			},
			storageReturn: 42,
			storageError:  nil,
			wantId:        42,
			wantError:     nil,
		},
		{
			name: "storage error",
			request: dto.PostTaskRequest{
				Name:        "Test Task",
				Description: "Test Description",
				Status:      model.Created,
			},
			storageReturn: -1,
			storageError:  errors.New("storage error"),
			wantId:        -1,
			wantError:     fmt.Errorf("TaskUsecase: couldn't store the task: storage error"),
		},
		{
			name: "invalid task status",
			request: dto.PostTaskRequest{
				Name:        "Test Task",
				Description: "Test Description",
				Status:      "invalid-status",
			},
			storageReturn: -1,
			storageError:  nil,
			wantId:        -1,
			wantError:     fmt.Errorf("TaskUsecase: couldn't store the task: invalid status in task: unknown type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			mockStorage := &MockTaskStorage{
				storeFunc: func(ctx context.Context, task model.Task) (int, error) {
					if tt.request.Status != task.Status ||
						tt.request.Name != task.Name ||
						tt.request.Description != task.Description {
						t.Error("Task conversion mismatch")
					}
					return tt.storageReturn, tt.storageError
				},
			}

			usecase, _ := NewTaskUsecase(mockLogger, mockStorage)
			gotId, gotErr := usecase.Store(ctx, tt.request)

			if gotId != tt.wantId {
				t.Errorf("Expected ID %d, got %d", tt.wantId, gotId)
			}

			if (gotErr == nil) != (tt.wantError == nil) {
				t.Fatalf("Error mismatch. Expected %v, got %v", tt.wantError, gotErr)
			}

			if gotErr != nil && gotErr.Error() != tt.wantError.Error() {
				t.Errorf("Error message mismatch. Expected %q, got %q", tt.wantError.Error(), gotErr.Error())
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	testTasks := []model.Task{
		{Id: 1, Name: "Task 1", Description: "Desc 1", Status: model.Created, CreatedAt: now},
		{Id: 2, Name: "Task 2", Description: "Desc 2", Status: model.InProgress, CreatedAt: now},
	}

	tests := []struct {
		name         string
		filter       model.Filter
		storageTasks []model.Task
		storageError error
		wantResponse dto.GetAllTasksResponse
		wantError    error
	}{
		{
			name:         "successful get all",
			filter:       model.Filter{Status: ""},
			storageTasks: testTasks,
			storageError: nil,
			wantResponse: dto.GetAllTasksResponse{
				Amount: 2,
				Tasks:  testTasks,
			},
			wantError: nil,
		},
		{
			name:         "filtered get all",
			filter:       model.Filter{Status: model.Created},
			storageTasks: []model.Task{testTasks[0]},
			storageError: nil,
			wantResponse: dto.GetAllTasksResponse{
				Amount: 1,
				Tasks:  []model.Task{testTasks[0]},
			},
			wantError: nil,
		},
		{
			name:         "storage error",
			filter:       model.Filter{Status: ""},
			storageTasks: nil,
			storageError: errors.New("storage error"),
			wantResponse: dto.GetAllTasksResponse{},
			wantError:    fmt.Errorf("TaskUsecase: couldn't get all the tasks: storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			mockStorage := &MockTaskStorage{
				getAllFunc: func(ctx context.Context, filter model.Filter) ([]model.Task, error) {
					if filter.Status != tt.filter.Status {
						t.Errorf("Filter mismatch. Expected %q, got %q", tt.filter.Status, filter.Status)
					}
					return tt.storageTasks, tt.storageError
				},
			}

			usecase, _ := NewTaskUsecase(mockLogger, mockStorage)
			gotResponse, gotErr := usecase.GetAll(ctx, tt.filter)

			if gotResponse.Amount != tt.wantResponse.Amount {
				t.Errorf("Amount mismatch. Expected %d, got %d", tt.wantResponse.Amount, gotResponse.Amount)
			}

			if len(gotResponse.Tasks) != len(tt.wantResponse.Tasks) {
				t.Fatalf("Tasks length mismatch. Expected %d, got %d", len(tt.wantResponse.Tasks), len(gotResponse.Tasks))
			}

			for i := range gotResponse.Tasks {
				if gotResponse.Tasks[i].Id != tt.wantResponse.Tasks[i].Id {
					t.Errorf("Task ID mismatch at index %d. Expected %d, got %d", i, tt.wantResponse.Tasks[i].Id, gotResponse.Tasks[i].Id)
				}
			}

			if (gotErr == nil) != (tt.wantError == nil) {
				t.Fatalf("Error mismatch. Expected %v, got %v", tt.wantError, gotErr)
			}

			if gotErr != nil && gotErr.Error() != tt.wantError.Error() {
				t.Errorf("Error message mismatch. Expected %q, got %q", tt.wantError.Error(), gotErr.Error())
			}
		})
	}
}

func TestGetByTaskId(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	testTask := model.Task{
		Id:          42,
		Name:        "Test Task",
		Description: "Test Description",
		Status:      model.Done,
		CreatedAt:   now,
	}

	tests := []struct {
		name         string
		taskId       int
		storageTask  *model.Task
		storageError error
		wantResponse dto.GetTaskByIdResponse
		wantError    error
	}{
		{
			name:         "successful get by id",
			taskId:       42,
			storageTask:  &testTask,
			storageError: nil,
			wantResponse: dto.GetTaskByIdResponse{
				Id:          42,
				Name:        "Test Task",
				Description: "Test Description",
				Status:      model.Done,
				CreatedAt:   now,
			},
			wantError: nil,
		},
		{
			name:         "task not found",
			taskId:       99,
			storageTask:  nil,
			storageError: fmt.Errorf("task not found"),
			wantResponse: dto.GetTaskByIdResponse{},
			wantError:    fmt.Errorf("TaskUsecase: task not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			mockStorage := &MockTaskStorage{
				getByTaskIdFunc: func(ctx context.Context, taskId int) (*model.Task, error) {
					if taskId != tt.taskId {
						t.Errorf("Task ID mismatch. Expected %d, got %d", tt.taskId, taskId)
					}
					return tt.storageTask, tt.storageError
				},
			}

			usecase, _ := NewTaskUsecase(mockLogger, mockStorage)
			gotResponse, gotErr := usecase.GetByTaskId(ctx, tt.taskId)

			if gotResponse.Id != tt.wantResponse.Id ||
				gotResponse.Name != tt.wantResponse.Name ||
				gotResponse.Description != tt.wantResponse.Description ||
				gotResponse.Status != tt.wantResponse.Status ||
				!gotResponse.CreatedAt.Equal(tt.wantResponse.CreatedAt) {
				t.Errorf("Response mismatch. Expected %v, got %v", tt.wantResponse, gotResponse)
			}

			if (gotErr == nil) != (tt.wantError == nil) {
				t.Fatalf("Error mismatch. Expected %v, got %v", tt.wantError, gotErr)
			}

			if gotErr != nil && gotErr.Error() != tt.wantError.Error() {
				t.Errorf("Error message mismatch. Expected %q, got %q", tt.wantError.Error(), gotErr.Error())
			}
		})
	}
}
