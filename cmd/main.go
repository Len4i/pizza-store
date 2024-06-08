package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"syscall"

	"github.com/Len4i/pizza-store/internal/config"
	mwLogger "github.com/Len4i/pizza-store/internal/middleware/logger"
	"github.com/Len4i/pizza-store/internal/services/order"
	"github.com/Len4i/pizza-store/internal/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.MustLoad()

	// Init logger
	logLevel := slog.Level(cfg.LogLevel)
	appLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	})
	log := slog.New(appLogHandler)
	// httpLog used in middleware for http access logs
	httpLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	httpLog := slog.New(httpLogHandler)

	log.Info(
		"starting pizza store api",
	)
	log.Debug("debug messages are enabled")

	// Init storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		os.Exit(1)
	}
	// Init services
	orderService := order.New(storage, log)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(mwLogger.New(httpLog))
	r.Use(middleware.Recoverer)

	r.Post("/order", orderService.Create)
	r.Get("/order/{id}", orderService.Get)
	// TODO: add DELETE /url/{id}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Non blocking start of the server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("failed to start server", "error", err)
			}
		}
	}()
	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", "error", err)
		return
	}

	// TODO: close storage connection

	log.Info("server stopped")
}
