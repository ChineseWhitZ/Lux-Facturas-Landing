package httpserver

import (
	"log/slog"
	"net/http"

	"lux-facturas/backend/internal/config"
)

type Server struct {
	config config.Config
	logger *slog.Logger
}

func NewServer(cfg config.Config, logger *slog.Logger) *Server {
	return &Server{
		config: cfg,
		logger: logger,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /api/v1/status", s.status)
	mux.HandleFunc("GET /openapi.yaml", s.openapi)
	mux.HandleFunc("GET /swagger", s.swagger)

	return requestLogger(s.logger, recoverer(cors(s.config.CORSOrigin, mux)))
}
