package models

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	TaskID      uuid.UUID   `json:"task_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     time.Time   `json:"end_time"`
	IsDone      bool        `json:"is_done"`
	DeletedAt   interface{} `json:"deleted_at"`
	UserId      uuid.UUID   `json:"user_id"`
}
