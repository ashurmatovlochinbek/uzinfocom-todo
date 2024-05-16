package user

import "net/http"

type Handler interface {
	Register() http.HandlerFunc
	Login() http.HandlerFunc
	CreateTask() http.HandlerFunc
	GetAllTasks() http.HandlerFunc
	DeleteTask() http.HandlerFunc
	UpdateTask() http.HandlerFunc
}
