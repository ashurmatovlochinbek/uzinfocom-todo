package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"uzinfocom-todo/internal/models"
	"uzinfocom-todo/internal/models/response_objects"
	"uzinfocom-todo/internal/user"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) user.Repository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user models.User) error {
	newUserUUID := uuid.New()
	createdUser := models.User{}
	createUserQuery := `INSERT INTO users (user_id, name, phone_number) VALUES ($1, $2, $3) RETURNING *`

	if err := r.db.QueryRowxContext(
		ctx,
		createUserQuery,
		&newUserUUID,
		&user.Name,
		&user.PhoneNumber,
	).Scan(&createdUser.UserId, &createdUser.Name, &createdUser.PhoneNumber); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByPhoneNumber(ctx context.Context, phone string) (*models.User, error) {
	query := `SELECT * FROM users WHERE phone_number = $1`
	user := models.User{}

	if err := r.db.QueryRowxContext(ctx, query, phone).Scan(&user.UserId, &user.Name, &user.PhoneNumber); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetTaskById(ctx context.Context, taskId uuid.UUID) (*models.Task, error) {
	query := `SELECT * FROM tasks WHERE task_id = $1`
	task := models.Task{}

	if err := r.db.QueryRowxContext(ctx, query, taskId).Scan(
		&task.TaskID,
		&task.Title,
		&task.Description,
		&task.StartTime,
		&task.EndTime,
		&task.IsDone,
		&task.DeletedAt,
		&task.UserId,
	); err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *userRepository) GetAllTasks(ctx context.Context) (*[]response_objects.ResponseTask, error) {
	userId := ctx.Value("user").(models.User).UserId
	query := `SELECT * FROM tasks WHERE user_id = $1 AND isDone IS FALSE AND deletedAt IS NULL`
	tasks := []response_objects.ResponseTask{}

	rows, err := r.db.QueryxContext(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := models.Task{}
		if err = rows.Scan(
			&t.TaskID,
			&t.Title,
			&t.Description,
			&t.StartTime,
			&t.EndTime,
			&t.IsDone,
			&t.DeletedAt,
			&t.UserId,
		); err != nil {
			return nil, err
		}

		task := response_objects.ResponseTask{
			TaskID:      t.TaskID,
			Title:       t.Title,
			Description: t.Description,
			StartTime:   t.StartTime.Format("02-01-2006 15:04"),
			EndTime:     t.EndTime.Format("02-01-2006 15:04"),
			IsDone:      t.IsDone,
			DeletedAt:   t.DeletedAt,
		}

		tasks = append(tasks, task)
	}

	return &tasks, nil

}

func (r *userRepository) CreateTask(ctx context.Context, task models.Task) error {
	userId := ctx.Value("user").(models.User).UserId
	taskId := uuid.New()
	createdTask := models.Task{}
	query := `INSERT INTO tasks (task_id, title, description, start_time, end_time, isDone, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
	var intf interface{}
	if err := r.db.QueryRowxContext(
		ctx,
		query,
		&taskId,
		&task.Title,
		&task.Description,
		&task.StartTime,
		&task.EndTime,
		&task.IsDone,
		&userId,
	).Scan(&createdTask.TaskID,
		&createdTask.Title,
		&createdTask.Description,
		&createdTask.StartTime,
		&createdTask.EndTime,
		&createdTask.IsDone,
		&intf,
		&createdTask.UserId); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteTask(ctx context.Context, taskId uuid.UUID) (sql.Result, error) {
	query := `UPDATE tasks SET deletedAt = CURRENT_TIMESTAMP WHERE task_id = $1`
	res, err := r.db.ExecContext(ctx, query, taskId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *userRepository) CheckForDeleteTask(ctx context.Context, taskId uuid.UUID) (int, error) {
	query := `SELECT count(*) from tasks where task_id = $1 and  deletedat is null`
	var count int

	if err := r.db.QueryRowxContext(ctx, query, taskId).Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (r *userRepository) UpdateTask(ctx context.Context, task models.Task) (sql.Result, error) {
	res, err := r.db.ExecContext(
		ctx,
		`UPDATE tasks SET title=$1, description=$2, start_time=$3, end_time=$4, isDone=$5`,
		&task.Title,
		&task.Description,
		&task.StartTime,
		&task.EndTime,
		&task.IsDone,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}
