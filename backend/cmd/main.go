package main

import (
	"context"
	"log"
	"os"

	"github.com/I-Missha/enterprise-management/config"
	"github.com/I-Missha/enterprise-management/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()
	logger := log.New(os.Stdout, "enterprise-management: ", log.LstdFlags)

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	err = dbpool.Ping(context.Background())
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	logger.Println("Successfully connected to database")

	router := chi.NewRouter()

	application := app.NewApp(cfg, logger, router, dbpool)
	if err := application.Run(); err != nil {
		logger.Fatalf("failed to run app: %v", err)
	}
}
