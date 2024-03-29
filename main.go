package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/config"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/repositories"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "InitConfig"))
	}

	db, err := initDb(ctx, cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "initDb"))
	}

	// запуск миграций на старте
	// упрощение - лучше делать через команды
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(errors.Wrap(err, "goose"))
	}
	if err := goose.Up(db.DB, "migrations"); err != nil {
		log.Fatal(errors.Wrap(err, "goose"))
	}

	repo := repositories.NewPgRepo(db)
	controller := api.NewController(repo)

	httpErrChan := listenHttp(initHttpServer(controller, repo, cfg))
	exitChan := startListenForQuit(ctx, cancel, cfg)
	log.Println("http server started at port " + strconv.Itoa(cfg.HttpPort))
	select {
	case err := <-httpErrChan:
		log.Fatal(errors.Wrap(err, "received error from http server"))
	case <-exitChan:
		break
	}
}

func initHttpServer(controller *api.Controller, service api.Service, cfg config.Config) *http.Server {
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
		Addr:    cfg.HttpHost + ":" + strconv.Itoa(cfg.HttpPort),
		Handler: r,
	}
}

func initDb(ctx context.Context, cfg config.Config) (*sqlx.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbSchema,
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

func startListenForQuit(ctx context.Context, ctxCancelFun context.CancelFunc, cfg config.Config) <-chan struct{} {
	exitChan := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 3)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return
		case sig := <-quit:
			log.Println("OS signal received: ", sig)
			ctxCancelFun()
			time.Sleep(time.Duration(cfg.CancelContextSleepDuration) * time.Second)
			exitChan <- struct{}{}
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
