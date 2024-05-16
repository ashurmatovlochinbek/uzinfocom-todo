package response_objects

import (
	"github.com/google/uuid"
)

type ResponseTask struct {
	TaskID      uuid.UUID   `json:"task_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	StartTime   string      `json:"start_time"`
	EndTime     string      `json:"end_time"`
	IsDone      bool        `json:"is_done"`
	DeletedAt   interface{} `json:"deleted_at"`
}
