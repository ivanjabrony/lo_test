package mapper

import (
	"ivanjabrony/test_lo/internal/model"
	"ivanjabrony/test_lo/internal/model/dto"
	"time"
)

func PostTaskRequestToTask(request dto.PostTaskRequest) model.Task {
	return model.Task{
		Id:          0,
		Status:      request.Status,
		Description: request.Description,
		Name:        request.Name,
		CreatedAt:   time.Time{}}
}

func TaskToGetTaskByIdReponse(task model.Task) dto.GetTaskByIdResponse {
	return dto.GetTaskByIdResponse{
		Id:          task.Id,
		Status:      task.Status,
		Description: task.Description,
		Name:        task.Name,
		CreatedAt:   task.CreatedAt}
}

func TasksToGetAllTasksResponse(tasks []model.Task) dto.GetAllTasksResponse {
	return dto.GetAllTasksResponse{
		Amount: len(tasks),
		Tasks:  tasks,
	}
}

func TaskToPostTaskReponse(task model.Task) dto.PostTaskResponse {
	return dto.PostTaskResponse{
		Id:        task.Id,
		CreatedAt: task.CreatedAt}
}
