package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mraiyuu/orders-api/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func (w http.ResponseWriter, r *http.Request)  {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", loadOrderRoutes)

	return router
}


func loadOrderRoutes(router chi.Router) {
	orderHandler := &handler.Order{}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{}", orderHandler.GetByID)
	router.Put("/", orderHandler.UpdateByID)
	router.Delete("/{}", orderHandler.DeleteByID)
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr: ":3000",
		Handler: a.router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server %w", err)
	}
	return nil 
}


