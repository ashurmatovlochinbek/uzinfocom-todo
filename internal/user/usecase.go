package user

import (
	"context"
	"github.com/google/uuid"
	"uzinfocom-todo/internal/models"
	"uzinfocom-todo/internal/models/response_objects"
	"uzinfocom-todo/pkg/http_errors"
)

type UseCase interface {
	Create(ctx context.Context, user models.User) http_errors.RestErr
	GetByPhoneNumber(ctx context.Context, phone string) (*models.User, http_errors.RestErr)
	GetAllTasks(ctx context.Context) (*[]response_objects.ResponseTask, http_errors.RestErr)
	CreateTask(ctx context.Context, task models.Task) http_errors.RestErr
	DeleteTask(ctx context.Context, taskId uuid.UUID) http_errors.RestErr
	GetTaskById(ctx context.Context, taskId uuid.UUID) (*models.Task, http_errors.RestErr)
	UpdateTask(ctx context.Context, task models.Task) http_errors.RestErr
}
