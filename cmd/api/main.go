package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marcusolivares/lending-backend-quick/internal/config"
	"github.com/marcusolivares/lending-backend-quick/internal/httpapi"
	"github.com/marcusolivares/lending-backend-quick/internal/loans"
	"github.com/marcusolivares/lending-backend-quick/internal/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer func() {
		if err := postgres.ClosePool(context.Background(), pool); err != nil {
			log.Printf("close database pool: %v", err)
		}
	}()

	repository := postgres.NewLoanApplicationRepository(pool)
	service := loans.NewService(repository)
	handler := httpapi.NewRouter(service)

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}()

	log.Printf("starting API on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	fmt.Println("server stopped")
}
