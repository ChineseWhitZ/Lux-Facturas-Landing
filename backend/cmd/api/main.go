package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"lux-facturas/backend/internal/config"
	"lux-facturas/backend/internal/httpserver"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.Load()
	server := httpserver.NewServer(cfg, logger)

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
