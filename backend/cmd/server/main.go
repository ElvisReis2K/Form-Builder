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

	"github.com/ElvisReis2K/Form-Builder/backend/internal/auth"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/config"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/database"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/forms"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/httpserver"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/openapi"
)

func main() {
	command := "run"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "run":
		runServer()
	case "migrate":
		runMigrations()
	case "openapi", "--print-openapi":
		writeOpenAPI()
	default:
		log.Fatalf("unknown command %q\n%s", command, usage())
	}
}

func runServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	bootstrapCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db, err := database.Open(bootstrapCtx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg.SessionSecret, cfg.SessionTTL)
	authHandler := auth.NewHandler(authService, cfg.CookieSecure)

	formsRepo := forms.NewRepository(db)
	formsService := forms.NewService(formsRepo)
	formsHandler := forms.NewHandler(authService, formsService)

	server := httpserver.New(cfg, db, authHandler, formsHandler)

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server listening on http://%s", cfg.Address)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}
}

func runMigrations() {
	direction := "up"
	if len(os.Args) > 2 {
		direction = os.Args[2]
	}

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.Open(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch direction {
	case "up":
		err = database.MigrateUp(ctx, db)
	case "down":
		err = database.MigrateDown(ctx, db)
	default:
		log.Fatalf("unknown migrate direction %q\n%s", direction, usage())
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("migrate %s completed", direction)
}

func writeOpenAPI() {
	if err := openapi.Write(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func usage() string {
	return fmt.Sprintf(`Usage:
  %s run
  %s migrate up
  %s migrate down
  %s openapi`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
