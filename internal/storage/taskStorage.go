package storage

import (
	"context"
	"fmt"
	"ivanjabrony/test_lo/internal/model"
	"sync"
	"time"
)

const storageName = "TaskStorage"

type Logger interface {
	Log(format string, info ...any)
}

type TaskStorage struct {
	tasks     []model.Task
	idCounter int
	logger    Logger
	m         sync.RWMutex
}

func NewTaskStorage(logger Logger) (*TaskStorage, error) {
	logger.Log("Created %s successfully", storageName)

	return &TaskStorage{
		make([]model.Task, 0),
		0,
		logger,
		sync.RWMutex{},
	}, nil
}

func (st *TaskStorage) Store(ctx context.Context, task model.Task) (int, error) {
	st.m.Lock()
	defer st.m.Unlock()
	task.Id = len(st.tasks)
	task.CreatedAt = time.Now()
	st.tasks = append(st.tasks, task)
	st.idCounter++

	st.logger.Log("Stored task: %v sucsessfully", task)

	return task.Id, nil
}

func (st *TaskStorage) GetAll(ctx context.Context, filter model.Filter) ([]model.Task, error) {
	st.m.RLock()
	defer st.m.RUnlock()
	ans := make([]model.Task, 0)

	for _, task := range st.tasks {
		if filter.Status == "" || task.Status == filter.Status {
			ans = append(ans, task)
		}
	}

	return ans, nil
}

func (st *TaskStorage) GetByTaskId(ctx context.Context, TaskId int) (*model.Task, error) {
	st.m.RLock()
	defer st.m.RUnlock()
	if TaskId > st.idCounter || TaskId < 0 {
		return nil, fmt.Errorf("%v: error while retrieving task by id(%v): nonexistent id", storageName, TaskId)
	}
	ans := st.tasks[TaskId]

	return &ans, nil
}
