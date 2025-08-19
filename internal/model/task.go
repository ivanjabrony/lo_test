package model

import (
	"time"
)

type TaskStatus string

type Task struct {
	Id          int        `json:"id"`
	Status      TaskStatus `json:"status"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
}
