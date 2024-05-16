package http

import (
	"github.com/go-chi/chi/v5"
	"uzinfocom-todo/config"
	"uzinfocom-todo/internal/middleware"
	"uzinfocom-todo/internal/user"
)

func MapRoutes(router *chi.Mux, h user.Handler, cfg *config.Config) {
	router.Route("/api", func(r chi.Router) {
		r.Post("/register", h.Register())
		r.Post("/login", h.Login())
	})

	router.Route("/todo", func(r chi.Router) {
		r.Use(middleware.AuthJwtMiddleware(cfg.Server.JwtSecretKey))
		r.Post("/", h.CreateTask())
		r.Get("/", h.GetAllTasks())
		r.Delete("/{taskId}", h.DeleteTask())
		r.Put("/{taskId}", h.UpdateTask())
	})
}
