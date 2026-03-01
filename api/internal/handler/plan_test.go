package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/handler"
	"github.com/sibukixxx/travelist/api/internal/infra/clients"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

func setupPlanHandler(t *testing.T) *handler.PlanHandler {
	t.Helper()
	placesClient := clients.NewStubPlacesClient()
	llmClient := clients.NewStubLLMClient()
	itineraryRepo := repo.NewMemoryItineraryRepository()
	generator := usecase.NewPlanGenerator(placesClient, llmClient, itineraryRepo, clock.RealClock{})
	return handler.NewPlanHandler(generator)
}

func TestPlanHandlerGeneratePlan(t *testing.T) {
	t.Run("returns 200 with JSON response for valid POST request", func(t *testing.T) {
		h := setupPlanHandler(t)
		body := `{"destination":"京都","num_days":2,"start_date":"2026-04-01","preferences":{"interests":["culture"],"budget":"moderate","travel_style":"balanced"}}`
		req := httptest.NewRequest(http.MethodPost, "/api/plans", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.GeneratePlan(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		contentType := rec.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
		}

		var result map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to parse response JSON: %v", err)
		}

		if result["itinerary"] == nil {
			t.Error("response should contain 'itinerary' field")
		}
		if result["budget_summary"] == nil {
			t.Error("response should contain 'budget_summary' field")
		}
	})

	t.Run("returns 405 for GET request", func(t *testing.T) {
		h := setupPlanHandler(t)
		req := httptest.NewRequest(http.MethodGet, "/api/plans", nil)
		rec := httptest.NewRecorder()

		h.GeneratePlan(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
		}
	})

	t.Run("returns 400 for invalid JSON body", func(t *testing.T) {
		h := setupPlanHandler(t)
		req := httptest.NewRequest(http.MethodPost, "/api/plans", strings.NewReader("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.GeneratePlan(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}

		var result map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to parse error response: %v", err)
		}
		if result["error"] == "" {
			t.Error("error response should contain 'error' field")
		}
	})

	t.Run("returns 500 for invalid start_date", func(t *testing.T) {
		h := setupPlanHandler(t)
		body := `{"destination":"京都","num_days":1,"start_date":"not-a-date","preferences":{"budget":"moderate"}}`
		req := httptest.NewRequest(http.MethodPost, "/api/plans", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.GeneratePlan(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})
}
