package app

import (
	"log"
	"net/http"

	"github.com/I-Missha/enterprise-management/config"
	"github.com/I-Missha/enterprise-management/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg    *config.Config
	logger *log.Logger
	router chi.Router
	dbpool *pgxpool.Pool
}

func NewApp(cfg *config.Config, logger *log.Logger, router chi.Router, dbpool *pgxpool.Pool) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
		router: router,
		dbpool: dbpool,
	}
}

func (a *App) Run() error {
	h := handler.NewHandler(a.logger, a.dbpool)
	h.InitRoutes(a.router)

	a.logger.Printf("starting server on %s", a.cfg.HTTPServer.Address)

	srv := &http.Server{
		Addr:         a.cfg.HTTPServer.Address,
		Handler:      a.router,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
		IdleTimeout:  a.cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
