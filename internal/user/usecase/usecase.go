package usecase

import (
	"context"
	"github.com/google/uuid"
	"time"
	"uzinfocom-todo/config"
	"uzinfocom-todo/internal/models"
	"uzinfocom-todo/internal/models/response_objects"
	"uzinfocom-todo/internal/user"
	"uzinfocom-todo/pkg/http_errors"
)

type userUseCase struct {
	userRepo user.Repository
	cfg      *config.Config
}

func NewUserUseCase(userRepo user.Repository, cfg *config.Config) user.UseCase {
	return &userUseCase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (uc *userUseCase) Create(ctx context.Context, user models.User) http_errors.RestErr {
	err := uc.userRepo.Create(ctx, user)

	if err != nil {
		return http_errors.ParseErrors(err)
	}
	return nil
}

func (uc *userUseCase) GetByPhoneNumber(ctx context.Context, phone string) (*models.User, http_errors.RestErr) {
	user, err := uc.userRepo.GetByPhoneNumber(ctx, phone)

	if err != nil {
		return nil, http_errors.ParseErrors(err)
	}
	return user, nil
}

func (uc *userUseCase) GetAllTasks(ctx context.Context) (*[]response_objects.ResponseTask, http_errors.RestErr) {
	tasks, err := uc.userRepo.GetAllTasks(ctx)
	if err != nil {
		return nil, http_errors.ParseErrors(err)
	}
	return tasks, nil
}

func (uc *userUseCase) CreateTask(ctx context.Context, task models.Task) http_errors.RestErr {
	responseTasks, _ := uc.userRepo.GetAllTasks(ctx)

	if len(*responseTasks) > 0 {
		for _, responseTask := range *responseTasks {
			checkTime1 := task.StartTime
			checkTime2 := task.EndTime
			responseTaskStartTime, _ := time.Parse("02-01-2006 15:04", responseTask.StartTime)
			responseTaskEndTime, _ := time.Parse("02-01-2006 15:04", responseTask.EndTime)
			if (checkTime1.After(responseTaskStartTime) && checkTime1.Before(responseTaskEndTime)) || (checkTime2.After(responseTaskStartTime) && checkTime2.Before(
				responseTaskEndTime)) || (checkTime1.Equal(responseTaskStartTime) || checkTime1.Equal(
				responseTaskEndTime)) || (checkTime2.Equal(
				responseTaskStartTime) || checkTime2.Equal(
				responseTaskEndTime) || (checkTime1.Before(responseTaskStartTime) && checkTime2.After(responseTaskEndTime))) {
				return http_errors.TaskExistsBetweenGivenTime()
			}

		}
	}

	err := uc.userRepo.CreateTask(ctx, task)

	if err != nil {
		return http_errors.ParseErrors(err)
	}
	return nil
}

func (uc *userUseCase) DeleteTask(ctx context.Context, taskId uuid.UUID) http_errors.RestErr {
	count, err := uc.userRepo.CheckForDeleteTask(ctx, taskId)
	if err != nil {
		return http_errors.ParseErrors(err)
	}
	if count > 0 {
		_, err := uc.userRepo.DeleteTask(ctx, taskId)
		if err != nil {
			return http_errors.ParseErrors(err)
		}
		return nil
	}

	return http_errors.ObjectNotFoundToDelete()
}

func (uc *userUseCase) GetTaskById(ctx context.Context, taskId uuid.UUID) (*models.Task, http_errors.RestErr) {
	t, err := uc.userRepo.GetTaskById(ctx, taskId)

	if err != nil {
		return nil, http_errors.ParseErrors(err)
	}

	return t, nil
}

func (uc *userUseCase) UpdateTask(ctx context.Context, task models.Task) http_errors.RestErr {

	//if task.IsDone == true {
	//	now := time.Now()
	//
	//	if now.Before(task.StartTime) {
	//		return http_errors.UpdateIsDoneErr()
	//	}
	//}

	responseTasks, _ := uc.userRepo.GetAllTasks(ctx)

	if len(*responseTasks) > 0 {
		for _, responseTask := range *responseTasks {
			checkTime1 := task.StartTime
			checkTime2 := task.EndTime
			responseTaskStartTime, _ := time.Parse("02-01-2006 15:04", responseTask.StartTime)
			responseTaskEndTime, _ := time.Parse("02-01-2006 15:04", responseTask.EndTime)

			if (checkTime1.After(responseTaskStartTime) && checkTime1.Before(responseTaskEndTime)) || (checkTime2.After(responseTaskStartTime) && checkTime2.Before(
				responseTaskEndTime)) || (checkTime1.Equal(responseTaskStartTime) || checkTime1.Equal(
				responseTaskEndTime)) || (checkTime2.Equal(
				responseTaskStartTime) || checkTime2.Equal(
				responseTaskEndTime) || (checkTime1.Before(responseTaskStartTime) && checkTime2.After(responseTaskEndTime))) {
				return http_errors.TaskExistsBetweenGivenTime()
			}
		}
	}

	_, err := uc.userRepo.UpdateTask(ctx, task)

	if err != nil {
		return http_errors.ParseErrors(err)
	}

	return nil
}
