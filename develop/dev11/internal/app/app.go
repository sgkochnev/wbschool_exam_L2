package app

import (
	"calendar/internal/config"
	"calendar/internal/http_server/router"
	"calendar/internal/storage"
	httpserver "calendar/pkg/http_server"
	"calendar/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// Run application
func Run() {

	cfg := config.New()

	log := logger.New(slog.LevelDebug)

	store := storage.New(cfg.StroagePath)

	routes := router.New(store)

	server := httpserver.New(routes, httpserver.WithAddress(cfg.Address))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("interrupt", slog.String("signal", s.String()))
	case err := <-server.Notify():
		log.Info("server stopped", slog.String("error", err.Error()))
	}

	if err := server.Shutdown(); err != nil {
		log.Error("server shutdown error", slog.String("error", err.Error()))
	}
	log.Info("server stopped")
}
