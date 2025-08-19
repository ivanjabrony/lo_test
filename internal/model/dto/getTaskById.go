package dto

import (
	"ivanjabrony/test_lo/internal/model"
	"time"
)

type GetTaskByIdResponse struct {
	Id          int              `json:"id"`
	Status      model.TaskStatus `json:"status"`
	Name        string           `json:"name"`
	Description string           `json:"amount"`
	CreatedAt   time.Time        `json:"created_at"`
}
