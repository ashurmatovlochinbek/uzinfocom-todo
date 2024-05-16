package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uzinfocom-todo/config"
	uhttp "uzinfocom-todo/internal/user/delivery/http"
	"uzinfocom-todo/internal/user/repository"
	"uzinfocom-todo/internal/user/usecase"
	"uzinfocom-todo/pkg/db/db_postgres"
)

func main() {
	log.Println("Starting api server")
	cfgFileViper, err := config.LoadConfig("./config/config-local")

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	cfg, err := config.ParseConfig(cfgFileViper)

	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}

	db, err := db_postgres.NewPsqlDB(cfg)

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("Error closing database: ", err)
		}
	}()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	if err != nil {
		log.Fatalf("PostgreSQLInstanceError: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)

	if err != nil {
		log.Fatalf("NewWithDatabaseInstanceError: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate.Up: %v", err)
	} else {
		log.Println("migrate.Up success")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	repo := repository.NewUserRepository(db)
	uc := usecase.NewUserUseCase(repo, cfg)
	uh := uhttp.NewUserHandler(cfg, uc)

	uhttp.MapRoutes(r, uh, cfg)

	server := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v", err)
	}

	log.Println("Server exiting gracefully")

	if err = m.Down(); err != nil {
		log.Fatalf("Down err: %v", err)
	} else {
		log.Println("Down success")
	}
}
