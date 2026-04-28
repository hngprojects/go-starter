package server

import (
	"log/slog"
	"net/http"

	"github.com/user/go-starter/internal/config"
	"github.com/user/go-starter/internal/handlers"
	"github.com/user/go-starter/internal/middleware"
)

type Server struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Server {
	return &Server{cfg: cfg, log: log}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	health := handlers.NewHealthHandler()
	mux.HandleFunc("GET /health", health.Check)
	mux.HandleFunc("GET /{$}", health.Root)

	chain := middleware.Chain(
		middleware.Recover(s.log),
		middleware.Logger(s.log),
	)
	return chain(mux)
}
