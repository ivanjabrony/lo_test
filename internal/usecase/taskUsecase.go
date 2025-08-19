package usecase

import (
	"context"
	"fmt"
	"ivanjabrony/test_lo/internal/model"
	"ivanjabrony/test_lo/internal/model/dto"
	"ivanjabrony/test_lo/internal/model/mapper"
)

const usecaseName = "TaskUsecase"

type TaskStorage interface {
	Store(ctx context.Context, task model.Task) (int, error)
	GetAll(ctx context.Context, filter model.Filter) ([]model.Task, error)
	GetByTaskId(ctx context.Context, taskId int) (*model.Task, error)
}

type Logger interface {
	Log(format string, info ...any)
}

type TaskUsecase struct {
	logger      Logger
	taskStorage TaskStorage
}

func NewTaskUsecase(logger Logger, storage TaskStorage) (*TaskUsecase, error) {
	if storage == nil {
		return nil, fmt.Errorf("nil values in %v constructor", usecaseName)
	}

	logger.Log("Created %s successfully", usecaseName)
	return &TaskUsecase{taskStorage: storage}, nil
}

func (tu *TaskUsecase) Store(ctx context.Context, request dto.PostTaskRequest) (int, error) {
	task := mapper.PostTaskRequestToTask(request)
	if err := model.ValidateTask(task); err != nil {
		return -1, fmt.Errorf("%v: couldn't store the task: %w", usecaseName, err)
	}
	id, err := tu.taskStorage.Store(ctx, task)
	if err != nil {
		return -1, fmt.Errorf("%v: couldn't store the task: %w", usecaseName, err)
	}
	return id, nil
}

func (tu *TaskUsecase) GetAll(ctx context.Context, filter model.Filter) (dto.GetAllTasksResponse, error) {
	tasks, err := tu.taskStorage.GetAll(ctx, filter)
	if err != nil {
		return dto.GetAllTasksResponse{}, fmt.Errorf("%v: couldn't get all the tasks: %w", usecaseName, err)
	}

	return mapper.TasksToGetAllTasksResponse(tasks), nil
}

func (tu *TaskUsecase) GetByTaskId(ctx context.Context, taskId int) (dto.GetTaskByIdResponse, error) {
	task, err := tu.taskStorage.GetByTaskId(ctx, taskId)
	if err != nil {
		return dto.GetTaskByIdResponse{}, fmt.Errorf("%v: %w", usecaseName, err)
	}
	response := mapper.TaskToGetTaskByIdReponse(*task)

	return response, err
}
