package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"lux-facturas/backend/internal/config"
	"lux-facturas/backend/internal/httpserver"
	"lux-facturas/backend/internal/leads"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.Load()
	leadRepository, err := leads.NewRepository(cfg.LeadsFile)
	if err != nil {
		logger.Error("could not initialize lead repository", "error", err)
		os.Exit(1)
	}

	server := httpserver.NewServer(cfg, logger, leadRepository)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      server.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("starting Lux Facturas API", "port", cfg.Port, "env", cfg.Environment)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("api server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
