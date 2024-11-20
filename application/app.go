package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		rdb:    redis.NewClient(&redis.Options{}),
	}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		fmt.Println("failed to connect to redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("starting server ")

	ch := make(chan error, 1)

	err = server.ListenAndServe()
	if err != nil {
		ch <- fmt.Errorf("failed to start server %w", err)
	}
	close(ch)

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancle := context.WithTimeout(context.Background(), time.Second*10)
		defer cancle()
		return server.Shutdown(timeout)
	}

	// return nil
}
