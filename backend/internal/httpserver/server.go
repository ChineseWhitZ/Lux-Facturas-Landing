package httpserver

import (
	"log/slog"
	"net/http"

	"lux-facturas/backend/internal/config"
	"lux-facturas/backend/internal/leads"
)

type Server struct {
	config         config.Config
	logger         *slog.Logger
	leadRepository *leads.Repository
}

func NewServer(cfg config.Config, logger *slog.Logger, leadRepository *leads.Repository) *Server {
	return &Server{
		config:         cfg,
		logger:         logger,
		leadRepository: leadRepository,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.home)
	mux.HandleFunc("GET /login", s.login)
	mux.HandleFunc("POST /login", s.loginSubmit)
	mux.HandleFunc("POST /logout", s.logout)
	mux.HandleFunc("GET /dashboard", s.dashboard)
	mux.HandleFunc("POST /dashboard/leads/{id}/status", s.updateLeadStatusFromDashboard)
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /api/v1/status", s.status)
	mux.HandleFunc("POST /api/v1/leads", s.createLead)
	mux.HandleFunc("GET /api/v1/leads", s.listLeads)
	mux.HandleFunc("GET /api/v1/leads/{id}", s.getLead)
	mux.HandleFunc("PATCH /api/v1/leads/{id}/status", s.updateLeadStatus)
	mux.HandleFunc("GET /openapi.yaml", s.openapi)
	mux.HandleFunc("GET /swagger", s.swagger)

	return requestLogger(s.logger, recoverer(cors(s.config.CORSOrigin, mux)))
}
