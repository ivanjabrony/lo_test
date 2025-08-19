package storage

import (
	"context"
	"fmt"
	"ivanjabrony/test_lo/internal/model"
	"testing"
)

// MockLogger is a mock implementation of the Logger interface for testing
type MockLogger struct {
	logs []string
}

func (m *MockLogger) Log(format string, info ...any) {
	m.logs = append(m.logs, fmt.Sprintf(format, info...))
}

func TestNewTaskStorage(t *testing.T) {
	mockLogger := &MockLogger{}
	storage, err := NewTaskStorage(mockLogger)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if storage == nil {
		t.Fatal("Expected storage to be initialized, got nil")
	}

	if len(mockLogger.logs) != 1 {
		t.Fatalf("Expected 1 log message, got %d", len(mockLogger.logs))
	}

	expectedLog := "Created TaskStorage successfully"
	if mockLogger.logs[0] != expectedLog {
		t.Errorf("Expected log message '%s', got '%s'", expectedLog, mockLogger.logs[0])
	}
}

func TestStore(t *testing.T) {
	mockLogger := &MockLogger{}
	storage, _ := NewTaskStorage(mockLogger)
	ctx := context.Background()

	task := model.Task{
		Name:        "Test Task",
		Description: "Test Description",
		Status:      model.Created,
	}

	t.Run("successful storage", func(t *testing.T) {
		id, err := storage.Store(ctx, task)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if id != 0 {
			t.Errorf("Expected ID 0 for first task, got %d", id)
		}

		if len(storage.tasks) != 1 {
			t.Fatalf("Expected 1 task in storage, got %d", len(storage.tasks))
		}

		storedTask := storage.tasks[0]
		if storedTask.Name != task.Name || storedTask.Description != task.Description || storedTask.Status != task.Status {
			t.Error("Stored task doesn't match input task")
		}

		if storedTask.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set, got zero time")
		}

		if len(mockLogger.logs) < 2 {
			t.Fatalf("Expected at least 2 log messages, got %d", len(mockLogger.logs))
		}
	})

	t.Run("concurrent storage", func(t *testing.T) {
		storage, _ := NewTaskStorage(mockLogger)
		tasksToAdd := 100
		results := make(chan int, tasksToAdd)

		for i := range tasksToAdd {
			go func(i int) {
				newTask := model.Task{
					Name:        fmt.Sprintf("Task %d", i),
					Description: fmt.Sprintf("Description %d", i),
					Status:      model.InProgress,
				}
				id, _ := storage.Store(ctx, newTask)
				results <- id
			}(i)
		}

		ids := make([]int, tasksToAdd)
		for i := range tasksToAdd {
			ids[i] = <-results
		}

		uniqueIDs := make(map[int]bool)
		for _, id := range ids {
			if uniqueIDs[id] {
				t.Errorf("Duplicate ID found: %d", id)
			}
			uniqueIDs[id] = true
		}

		if len(storage.tasks) != tasksToAdd {
			t.Errorf("Expected %d tasks in storage, got %d", tasksToAdd, len(storage.tasks))
		}
	})
}

func TestGetByTaskId(t *testing.T) {
	mockLogger := &MockLogger{}
	storage, _ := NewTaskStorage(mockLogger)
	ctx := context.Background()

	tasks := []model.Task{
		{Name: "Task 1", Description: "Desc 1", Status: model.Created},
		{Name: "Task 2", Description: "Desc 2", Status: model.InProgress},
		{Name: "Task 3", Description: "Desc 3", Status: model.Done},
	}

	for i := range tasks {
		_, err := storage.Store(ctx, tasks[i])
		if err != nil {
			t.Fatalf("Failed to setup test: %v", err)
		}
	}

	tests := []struct {
		name     string
		taskId   int
		wantTask *model.Task
		wantErr  error
	}{
		{
			name:     "existing task",
			taskId:   1,
			wantTask: &tasks[1],
			wantErr:  nil,
		},
		{
			name:     "non-existent task",
			taskId:   99,
			wantTask: nil,
			wantErr:  fmt.Errorf("%v: error while retrieving task by id(99): nonexistent id", storageName),
		},
		{
			name:     "negative id",
			taskId:   -1,
			wantTask: nil,
			wantErr:  fmt.Errorf("%v: error while retrieving task by id(-1): nonexistent id", storageName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTask, gotErr := storage.GetByTaskId(ctx, tt.taskId)

			if tt.wantTask == nil && gotTask != nil {
				t.Errorf("Expected nil task, got %v", gotTask)
			} else if tt.wantTask != nil && gotTask == nil {
				t.Errorf("Expected task %v, got nil", tt.wantTask)
			} else if tt.wantTask != nil && gotTask != nil {
				if gotTask.Name != tt.wantTask.Name || gotTask.Description != tt.wantTask.Description || gotTask.Status != tt.wantTask.Status {
					t.Errorf("Task mismatch. Expected %v, got %v", tt.wantTask, gotTask)
				}
			}

			if (gotErr == nil) != (tt.wantErr == nil) {
				t.Errorf("Error mismatch. Expected %v, got %v", tt.wantErr, gotErr)
			} else if gotErr != nil && gotErr.Error() != tt.wantErr.Error() {
				t.Errorf("Error message mismatch. Expected %q, got %q", tt.wantErr.Error(), gotErr.Error())
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	mockLogger := &MockLogger{}
	storage, _ := NewTaskStorage(mockLogger)
	ctx := context.Background()

	tasks := []model.Task{
		{Name: "Task 1", Description: "Desc 1", Status: model.Created},
		{Name: "Task 2", Description: "Desc 2", Status: model.InProgress},
		{Name: "Task 3", Description: "Desc 3", Status: model.Done},
		{Name: "Task 4", Description: "Desc 4", Status: model.Created},
	}

	for i := range tasks {
		_, err := storage.Store(ctx, tasks[i])
		if err != nil {
			t.Fatalf("Failed to setup test: %v", err)
		}
	}

	tests := []struct {
		name    string
		filter  model.Filter
		wantLen int
	}{
		{
			name:    "no filter",
			filter:  model.Filter{Status: ""},
			wantLen: 4,
		},
		{
			name:    "filter by todo",
			filter:  model.Filter{Status: model.Created},
			wantLen: 2,
		},
		{
			name:    "filter by in progress",
			filter:  model.Filter{Status: model.InProgress},
			wantLen: 1,
		},
		{
			name:    "filter by done",
			filter:  model.Filter{Status: model.Done},
			wantLen: 1,
		},
		{
			name:    "filter by non-existent status",
			filter:  model.Filter{Status: "non-existent"},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTasks, err := storage.GetAll(ctx, tt.filter)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(gotTasks) != tt.wantLen {
				t.Errorf("Expected %d tasks, got %d", tt.wantLen, len(gotTasks))
			}

			for _, task := range gotTasks {
				if tt.filter.Status != "" && task.Status != tt.filter.Status {
					t.Errorf("Expected all tasks to have status %q, got %q", tt.filter.Status, task.Status)
				}
			}
		})
	}
}
