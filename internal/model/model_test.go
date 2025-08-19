package model

import (
	"testing"
	"time"
)

func TestTaskValidation(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr bool
	}{
		{
			name: "valid done task",
			task: Task{
				Id:          1,
				Status:      Done,
				Name:        "task 1",
				Description: "descr 1",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid in progress task",
			task: Task{
				Id:          1,
				Status:      InProgress,
				Name:        "task 1",
				Description: "descr 1",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid created task",
			task: Task{
				Id:          1,
				Status:      Created,
				Name:        "task 1",
				Description: "descr 1",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid status task",
			task: Task{
				Id:          1,
				Status:      "invalid",
				Name:        "task 1",
				Description: "descr 1",
				CreatedAt:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "negative id",
			task: Task{
				Id:          -1,
				Status:      "invalid",
				Name:        "task 1",
				Description: "descr 1",
				CreatedAt:   time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTask(tt.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilterValidation(t *testing.T) {
	tests := []struct {
		name    string
		filter  Filter
		wantErr bool
	}{
		{
			name:    "valid done filter",
			filter:  Filter{Status: Done},
			wantErr: false,
		},
		{
			name:    "valid in progress filter",
			filter:  Filter{Status: InProgress},
			wantErr: false,
		},
		{
			name:    "valid created filter",
			filter:  Filter{Status: Created},
			wantErr: false,
		},
		{
			name:    "invalid status filter",
			filter:  Filter{Status: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilter(tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
