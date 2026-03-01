package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sibukixxx/travelist/api/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/health", handler.HealthCheck)

	// Serve frontend static files (production mode)
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir != "" {
		fs := http.FileServer(http.Dir(staticDir))
		mux.Handle("/", fs)
		log.Printf("Serving static files from %s", staticDir)
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
