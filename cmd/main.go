package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RakibulBh/AI-pr-reviewer/internal/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	config.Bootstrap(&config.BootstrapConfig{
		R: r,
	})

	fmt.Printf("test")

	// start the server
	log.Fatal(http.ListenAndServe(":3000", r))
}
