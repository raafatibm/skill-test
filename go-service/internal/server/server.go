package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reporting/internal/handlers"
	"reporting/internal/models"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	cfg models.Config

	httpServer *http.Server
}

func NewServer(cfg models.Config) (*Server, error) {
	server := Server{
		cfg: cfg,
	}

	handler, err := handlers.NewHandler(cfg)
	if err != nil {
		return nil, err
	}

	router := setRouter(handler)

	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.ListenPort),
		Handler: router,
	}

	return &server, nil
}

func (srv *Server) Start() {
	log.Println("http server has been started at port ", srv.cfg.ListenPort)
	err := srv.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Println("http server stopped unexpected")
	} else {
		log.Println("http server process stopped")
	}
}

func (srv *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.httpServer.Shutdown(ctx)
}

func setRouter(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/api/v1/students/{student_id}/report", handler.GetStudentDetails)

	return r
}
