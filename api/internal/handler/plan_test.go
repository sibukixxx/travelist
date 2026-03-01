package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"encoding/json"

	"github.com/sibukixxx/travelist/api/internal/handler"
	"github.com/sibukixxx/travelist/api/internal/infra/clients"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

func setupHandler(t *testing.T) *handler.PlanHandler {
	t.Helper()
	placesClient := clients.NewStubPlacesClient()
	llmClient := clients.NewStubLLMClient()
	itineraryRepo := repo.NewInMemoryItineraryRepository()
	clk := clock.FixedClock{Time: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)}
	generator := usecase.NewPlanGenerator(placesClient, llmClient, itineraryRepo, clk)
	return handler.NewPlanHandler(generator)
}

func TestGeneratePlan(t *testing.T) {
	t.Run("POST /api/plans returns 200 with itinerary JSON", func(t *testing.T) {
		h := setupHandler(t)

		body := `{"destination":"京都","num_days":2,"start_date":"2026-04-01","preferences":{"interests":["culture"],"budget":"moderate","travel_style":"balanced"}}`
		req := httptest.NewRequest(http.MethodPost, "/api/plans", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.GeneratePlan(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		ct := resp.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}

		var result usecase.GenerateResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if result.Itinerary == nil {
			t.Fatal("expected itinerary in response, got nil")
		}
		if result.Itinerary.Destination != "京都" {
			t.Errorf("expected destination 京都, got %q", result.Itinerary.Destination)
		}
		if len(result.Itinerary.Days) != 2 {
			t.Errorf("expected 2 days, got %d", len(result.Itinerary.Days))
		}
		if result.BudgetSummary == nil {
			t.Error("expected budget_summary in response, got nil")
		}
	})

	t.Run("GET /api/plans returns 405 Method Not Allowed", func(t *testing.T) {
		h := setupHandler(t)

		req := httptest.NewRequest(http.MethodGet, "/api/plans", nil)
		w := httptest.NewRecorder()

		h.GeneratePlan(w, req)

		if w.Result().StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("expected 405, got %d", w.Result().StatusCode)
		}
	})

	t.Run("POST with invalid JSON returns 400", func(t *testing.T) {
		h := setupHandler(t)

		req := httptest.NewRequest(http.MethodPost, "/api/plans", strings.NewReader("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.GeneratePlan(w, req)

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Result().StatusCode)
		}
	})

	t.Run("POST with missing start_date returns 500", func(t *testing.T) {
		h := setupHandler(t)

		body := `{"destination":"京都","num_days":2,"preferences":{}}`
		req := httptest.NewRequest(http.MethodPost, "/api/plans", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.GeneratePlan(w, req)

		// start_date="" fails parsing → 500
		if w.Result().StatusCode != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Result().StatusCode)
		}
	})
}
