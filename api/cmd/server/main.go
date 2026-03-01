package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sibukixxx/travelist/api/internal/handler"
	"github.com/sibukixxx/travelist/api/internal/infra/clients"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Build dependencies
	placesClient := clients.NewStubPlacesClient()
	llmClient := clients.NewStubLLMClient()
	itineraryRepo := repo.NewInMemoryItineraryRepository()
	clk := clock.RealClock{}

	generator := usecase.NewPlanGenerator(placesClient, llmClient, itineraryRepo, clk)
	planHandler := handler.NewPlanHandler(generator)

	// API routes
	mux.HandleFunc("/api/health", handler.HealthCheck)
	mux.HandleFunc("/api/plans", planHandler.GeneratePlan)

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
