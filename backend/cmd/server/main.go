// Command server is the entrypoint of the Go version of the tobz-api application.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tobz/tobz-api/internal/config"
	"github.com/tobz/tobz-api/internal/database"
	"github.com/tobz/tobz-api/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi gagal: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database gagal: %v", err)
	}

	app := server.New(cfg, db)

	go func() {
		addr := ":" + cfg.Port
		log.Printf("tobz-api (Go) listening on %s [env=%s]", addr, cfg.AppEnv)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("server berhenti: %v", err)
		}
	}()

	// Graceful shutdown on termination signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	log.Println("bye 👋")
}
