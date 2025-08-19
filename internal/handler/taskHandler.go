package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"ivanjabrony/test_lo/internal/model"
	"ivanjabrony/test_lo/internal/model/dto"
	"net/http"
	"strconv"
)

const handlerName = "TransactionHandler"

type TaskUsecase interface {
	Store(ctx context.Context, request dto.PostTaskRequest) (int, error)
	GetAll(ctx context.Context, filter model.Filter) (dto.GetAllTasksResponse, error)
	GetByTaskId(ctx context.Context, taskId int) (dto.GetTaskByIdResponse, error)
}

type Logger interface {
	Log(format string, info ...any)
}

type TaskHandler struct {
	taskUsecase TaskUsecase
	logger      Logger
}

func NewTaskHandler(logger Logger, taskUsecase TaskUsecase) (*TaskHandler, error) {
	if taskUsecase == nil {
		return nil, fmt.Errorf("nil values in %v constructor", handlerName)
	}

	return &TaskHandler{taskUsecase, logger}, nil
}

func (th *TaskHandler) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var postReq dto.PostTaskRequest
	json.NewDecoder(r.Body).Decode(&postReq)

	unvalidatedTask := model.Task{
		Status:      postReq.Status,
		Name:        postReq.Name,
		Description: postReq.Description,
	}
	if err := model.ValidateTask(unvalidatedTask); err != nil {
		th.logger.Log("error in error %v: error while task validation: %v", handlerName, err)
		respondWithError(th.logger, w, http.StatusBadRequest, "invalid data in task")
		return
	}

	transactions, err := th.taskUsecase.Store(ctx, postReq)
	if err != nil {
		respondWithError(th.logger, w, http.StatusInternalServerError, "failed to store task")
		return
	}

	respondWithJSON(w, http.StatusOK, transactions)
}

func (th *TaskHandler) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queryParams := r.URL.Query()
	filter := model.EmptyFilter

	if typeParam := queryParams.Get("status"); typeParam != "" {
		unvalidatedFilter := model.Filter{Status: model.TaskStatus(typeParam)}
		if err := model.ValidateFilter(unvalidatedFilter); err != nil {
			th.logger.Log("error in error %v: error while filter validation: %v", handlerName, err)
			respondWithError(th.logger, w, http.StatusBadRequest, "invalid task status in filter")
			return
		}
		filter = unvalidatedFilter
	}

	response, err := th.taskUsecase.GetAll(ctx, filter)
	if err != nil {
		respondWithError(th.logger, w, http.StatusInternalServerError, "failed to retrieve tasks")
		th.logger.Log("error in %v: %v", handlerName, err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (th *TaskHandler) HandleGetTaskById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskIDParam := r.PathValue("task_id")
	if taskIDParam == "" {
		respondWithError(th.logger, w, http.StatusBadRequest, "task_id wasn't provided")
		return
	}

	taskId, err := strconv.Atoi(taskIDParam)
	if err != nil {
		respondWithError(th.logger, w, http.StatusBadRequest, "invalid task_id parameter")
		return
	}

	transactions, err := th.taskUsecase.GetByTaskId(ctx, taskId)
	if err != nil {
		respondWithError(th.logger, w, http.StatusInternalServerError, "failed to retrieve task")
		return
	}

	respondWithJSON(w, http.StatusOK, transactions)
}

func respondWithError(logger Logger, w http.ResponseWriter, code int, message string) {
	logger.Log("error in %v: %v", handlerName, message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
