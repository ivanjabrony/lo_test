package dto

import (
	"ivanjabrony/test_lo/internal/model"
	"time"
)

type PostTaskRequest struct {
	Status      model.TaskStatus `json:"status"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
}

type PostTaskResponse struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
