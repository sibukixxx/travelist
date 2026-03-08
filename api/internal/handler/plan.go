package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sibukixxx/travelist/api/internal/apperror"
	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

// PlanGeneratorInterface is the interface for plan generation.
type PlanGeneratorInterface interface {
	Generate(ctx context.Context, req domain.PlanRequest) (*usecase.GenerateResult, error)
}

// PlanHandler handles HTTP requests for itinerary planning.
type PlanHandler struct {
	generator PlanGeneratorInterface
}

// NewPlanHandler creates a new PlanHandler.
func NewPlanHandler(generator PlanGeneratorInterface) *PlanHandler {
	return &PlanHandler{generator: generator}
}

// GeneratePlan handles POST /api/plans.
func (h *PlanHandler) GeneratePlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.PlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, apperror.NewBadRequest("invalid request body"))
		return
	}

	// Apply defaults if not set
	if req.Constraint.MaxWalkDistanceM == 0 {
		req.Constraint = domain.DefaultConstraint()
	}

	result, err := h.generator.Generate(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// HealthCheck handles GET /api/health.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
