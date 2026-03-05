package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sibukixxx/travelist/api/internal/handler"
	"github.com/sibukixxx/travelist/api/internal/infra/clients"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/infra/email"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
	"github.com/sibukixxx/travelist/api/internal/infra/repo/sqlite"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// SQLite
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/travelist.db"
	}

	// Ensure data directory exists for file-based DB
	if dbPath != ":memory:" {
		if err := os.MkdirAll("./data", 0o755); err != nil {
			log.Fatalf("failed to create data directory: %v", err)
		}
	}

	db, err := sqlite.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Build dependencies
	clk := clock.RealClock{}

	// Plan generation dependencies
	placesClient := clients.NewStubPlacesClient()
	llmClient := clients.NewStubLLMClient()
	itineraryRepo := repo.NewMemoryItineraryRepository()
	planGenerator := usecase.NewPlanGenerator(placesClient, llmClient, itineraryRepo, clk)
	planHandler := handler.NewPlanHandler(planGenerator)

	// User registration dependencies
	userRepo := sqlite.NewUserRepo(db)
	emailSender := &email.LogEmailSender{}
	registrar := usecase.NewUserRegistrar(userRepo, emailSender, clk)
	verifier := usecase.NewEmailVerifier(userRepo, clk)
	userHandler := handler.NewUserHandler(registrar, verifier)

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/health", handler.HealthCheck)
	mux.HandleFunc("/api/plans", planHandler.GeneratePlan)
	mux.HandleFunc("/api/users", userHandler.Register)
	mux.HandleFunc("/api/users/verify", userHandler.VerifyEmail)

	// Serve frontend static files (production mode)
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir != "" {
		fs := http.FileServer(http.Dir(staticDir))
		mux.Handle("/", fs)
		log.Printf("Serving static files from %s", staticDir)
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s (db: %s)", addr, dbPath)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
