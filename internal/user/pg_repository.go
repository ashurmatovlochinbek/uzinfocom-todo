package user

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"uzinfocom-todo/internal/models"
	"uzinfocom-todo/internal/models/response_objects"
)

type Repository interface {
	Create(ctx context.Context, user models.User) error
	GetByPhoneNumber(ctx context.Context, phone string) (*models.User, error)
	CreateTask(ctx context.Context, task models.Task) error
	GetAllTasks(ctx context.Context) (*[]response_objects.ResponseTask, error)
	DeleteTask(ctx context.Context, taskId uuid.UUID) (sql.Result, error)
	CheckForDeleteTask(ctx context.Context, taskId uuid.UUID) (int, error)
	GetTaskById(ctx context.Context, taskId uuid.UUID) (*models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) (sql.Result, error)
}
