package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
	"uzinfocom-todo/config"
	"uzinfocom-todo/internal/models"
	"uzinfocom-todo/internal/models/request_objects"
	"uzinfocom-todo/internal/models/response_objects"
	"uzinfocom-todo/internal/user"
	"uzinfocom-todo/pkg/util"
)

type userHandler struct {
	cfg    *config.Config
	userUC user.UseCase
}

func NewUserHandler(cfg *config.Config, userUC user.UseCase) user.Handler {
	return &userHandler{
		cfg:    cfg,
		userUC: userUC,
	}
}

func (h *userHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var responseObject *response_objects.ResponseObject
		var u models.User

		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			responseObject = response_objects.NewResponseObject(false, "json format is incorrect for user", nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}

		defer r.Body.Close()

		restErr := h.userUC.Create(r.Context(), u)

		if restErr != nil {
			responseObject = response_objects.NewResponseObject(false, restErr.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}

		w.WriteHeader(http.StatusCreated)
		responseObject = response_objects.NewResponseObject(true, "user created", nil)
		jsonResponse, _ := json.Marshal(responseObject)
		w.Write(jsonResponse)
	}
}

func (h *userHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var responseObject *response_objects.ResponseObject
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			responseObject = response_objects.NewResponseObject(false, "json format is incorrect for user object", nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}

		defer r.Body.Close()

		u, restErr := h.userUC.GetByPhoneNumber(r.Context(), user.PhoneNumber)

		if restErr != nil {
			responseObject = response_objects.NewResponseObject(false, restErr.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(restErr.Status())
			w.Write(jsonResponse)
			return
		}

		accessToken, err := util.GenerateJWTToken(u, h.cfg)

		if err != nil {
			responseObject = response_objects.NewResponseObject(false, err.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonResponse)
			return
		}

		refreshToken, err := util.GenerateRefreshToken(u, h.cfg)

		if err != nil {
			responseObject = response_objects.NewResponseObject(false, err.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonResponse)
			return
		}

		responseObject = response_objects.NewResponseObject(true, "Access and Refresh Tokens generated successfully", map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
		jsonResponse, _ := json.Marshal(responseObject)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func (h *userHandler) GetAllTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var responseObject *response_objects.ResponseObject

		responseTasks, restErr := h.userUC.GetAllTasks(r.Context())

		if restErr != nil {
			responseObject = response_objects.NewResponseObject(false, restErr.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(restErr.Status())
			w.Write(jsonResponse)
			return
		}

		responseObject = response_objects.NewResponseObject(true, "Tasks are fetched", responseTasks)
		jsonResponse, _ := json.Marshal(responseObject)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func (h *userHandler) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var responseObject *response_objects.ResponseObject
		var requestTask request_objects.RequestTask

		err := json.NewDecoder(r.Body).Decode(&requestTask)

		if err != nil {
			responseObject = response_objects.NewResponseObject(false, err.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}

		defer r.Body.Close()

		startTime, _ := time.Parse("02-01-2006 15:04", requestTask.StartTime)
		endTime, _ := time.Parse("02-01-2006 15:04", requestTask.EndTime)

		task := models.Task{
			Title:       requestTask.Title,
			Description: requestTask.Description,
			StartTime:   startTime,
			EndTime:     endTime,
			IsDone:      false,
		}

		restErr := h.userUC.CreateTask(r.Context(), task)

		if restErr != nil {
			responseObject = response_objects.NewResponseObject(false, restErr.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(restErr.Status())
			w.Write(jsonResponse)
			return
		}

		responseObject = response_objects.NewResponseObject(true, "task created successfully", nil)
		jsonResponse, _ := json.Marshal(responseObject)
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}

func (h *userHandler) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var responseObject *response_objects.ResponseObject
		taskId, _ := uuid.Parse(chi.URLParam(r, "taskId"))

		restErr := h.userUC.DeleteTask(r.Context(), taskId)

		if restErr != nil {
			responseObject = response_objects.NewResponseObject(false, restErr.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(restErr.Status())
			w.Write(jsonResponse)
			return
		}

		responseObject = response_objects.NewResponseObject(true, "task deleted successfully", nil)
		jsonResponse, _ := json.Marshal(responseObject)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func (h *userHandler) UpdateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var responseObject *response_objects.ResponseObject
		var requestTask request_objects.RequestTaskForUpdate
		taskId, _ := uuid.Parse(chi.URLParam(r, "taskId"))
		err := json.NewDecoder(r.Body).Decode(&requestTask)
		var startTime time.Time
		var endTime time.Time

		if err != nil {
			responseObject = response_objects.NewResponseObject(false, err.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}

		defer r.Body.Close()

		t, restErr := h.userUC.GetTaskById(r.Context(), taskId)

		if restErr != nil {
			responseObject = response_objects.NewResponseObject(false, restErr.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(restErr.Status())
			w.Write(jsonResponse)
			return
		}

		if requestTask.Title != "" {
			t.Title = requestTask.Title
		}

		if requestTask.Description != "" {
			t.Description = requestTask.Description
		}

		if requestTask.IsDone == true {
			now := time.Now()

			if now.Before(t.StartTime) {
				responseObject = response_objects.NewResponseObject(true, "to set up isDone to true you have to be in the interval of task times", nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(jsonResponse)
				return
			}

			t.IsDone = true
		}

		if requestTask.StartTime != "" && requestTask.EndTime != "" {
			startTime, _ = time.Parse("02-01-2006 15:04", requestTask.StartTime)
			endTime, _ = time.Parse("02-01-2006 15:04", requestTask.EndTime)
			t.StartTime = startTime
			t.EndTime = endTime
		}

		restErr = h.userUC.UpdateTask(r.Context(), *t)

		if restErr != nil {
			responseObject = response_objects.NewResponseObject(false, restErr.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(restErr.Status())
			w.Write(jsonResponse)
			return
		}

		responseObject = response_objects.NewResponseObject(true, "task updated successfully", nil)
		jsonResponse, _ := json.Marshal(responseObject)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}
