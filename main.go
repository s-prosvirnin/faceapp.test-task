package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/repositories"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	cancelContextSleepSec = 2
	httpHost              = "localhost"
	httpPort              = "8085"

	dbHost     = "localhost"
	dbPort     = 54323
	dbUser     = "postgres"
	dbPassword = "12345"
	dbSchema   = "postgres"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-signalChan
		cancel()
		time.Sleep(cancelContextSleepSec * time.Second)
	}(ctx)

	db, err := initDb(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "initDb"))
	}
	repo := repositories.NewPgRepo(db)
	controller := api.NewController(repo)

	httpErrChan := listenHttp(initHttpServer(controller, repo))
	exitChan := startListenForQuit(ctx)
	select {
	case err := <-httpErrChan:
		log.Fatal(errors.Wrap(err, "received error from server"))
	case <-exitChan:
		break
	}
}

func initHttpServer(controller *api.Controller, service api.Service) *http.Server {
	r := mux.NewRouter()
	middleware := api.NewMiddleware(service)

	r.HandleFunc("/auth", controller.Auth).Methods("POST")
	r.Handle(
		"/team/contest",
		middleware.AuthRequest(controller.GetContest),
	).Methods("POST")
	r.Handle(
		"/team/contest/tasks",
		middleware.AuthRequest(controller.GetTeamTasks),
	).Methods("POST")
	r.Handle(
		"/team/contest/task/start",
		middleware.AuthRequest(controller.StartTask),
	).Methods("POST")
	r.Handle(
		"/team/contest/task/hint",
		middleware.AuthRequest(controller.ShowTaskHint),
	).Methods("POST")
	r.Handle(
		"/team/contest/task/answer",
		middleware.AuthRequest(controller.SendTaskAnswer),
	).Methods("POST")
	r.Handle(
		"/contest/results",
		middleware.AuthRequest(controller.GetContestResults),
	).Methods("POST")

	r.Use(middleware.MutateResponseHeaders)

	return &http.Server{
		Addr:    httpHost + ":" + httpPort,
		Handler: r,
	}
}

func initDb(ctx context.Context) (*sqlx.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbSchema,
	)

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, errors.Wrap(err, "creating postgres error")
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "postgres error")
	}

	return db, nil
}

func startListenForQuit(ctx context.Context) <-chan struct{} {
	exitChan := make(chan struct{})
	go func() {
		// @todo: pros не обрабатывается Ctrl+C
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return
		case <-quit:
			close(exitChan)

			return
		}
	}()

	return exitChan
}

func listenHttp(server *http.Server) <-chan error {
	errChan := make(chan error)

	go func() {
		errChan <- server.ListenAndServe()
	}()

	return errChan
}
