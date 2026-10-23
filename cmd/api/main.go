package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/therionlabs/lead-capture-backend/internal/auth"
	"github.com/therionlabs/lead-capture-backend/internal/catalogo"
	"github.com/therionlabs/lead-capture-backend/internal/config"
	"github.com/therionlabs/lead-capture-backend/internal/db"
	"github.com/therionlabs/lead-capture-backend/internal/prospecto"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error("el servidor terminó con error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// En desarrollo se puede usar un archivo .env; en otros entornos basta con las variables.
	if err := godotenv.Load(); err == nil {
		slog.Info("variables cargadas desde .env")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
		return err
	}
	slog.Info("migraciones al día")

	authRepo := auth.NewRepository(pool)
	authHandler := auth.NewHandler(authRepo)
	catalogoHandler := catalogo.NewHandler(catalogo.NewRepository(pool))
	prospectoHandler := prospecto.NewHandler(prospecto.NewRepository(pool), prospecto.NewValidator())

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger) // registra método, ruta, estado y duración; nunca el cuerpo
	r.Use(middleware.Recoverer)

	soloCaptador := auth.RequireRol(authRepo, auth.RolCaptador)
	soloAdministrativo := auth.RequireRol(authRepo, auth.RolAdministrativo)
	cualquierRol := auth.RequireRol(authRepo, auth.RolCaptador, auth.RolAdministrativo)

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)
		r.With(cualquierRol).Post("/auth/logout", authHandler.Logout)

		r.With(cualquierRol).Get("/catalogo", catalogoHandler.Get)
		r.With(soloCaptador).Post("/prospectos/sync", prospectoHandler.Sync)
		r.With(soloAdministrativo).Get("/prospectos", prospectoHandler.List)
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("API escuchando", "puerto", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("apagando servidor")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
