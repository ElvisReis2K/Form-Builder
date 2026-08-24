package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ElvisReis2K/Form-Builder/backend/internal/config"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/httpserver"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/openapi"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--print-openapi" {
		if err := openapi.Write(os.Stdout); err != nil {
			log.Fatal(err)
		}

		return
	}

	cfg := config.Load()
	server := httpserver.New(cfg)

	log.Printf("server listening on http://%s", cfg.Address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
