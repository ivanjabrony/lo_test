package dto

import (
	"ivanjabrony/test_lo/internal/model"
)

type GetAllTasksResponse struct {
	Amount int          `json:"amount"`
	Tasks  []model.Task `json:"tasks"`
}
