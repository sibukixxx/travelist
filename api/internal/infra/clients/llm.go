package clients

import "context"

// LLMClient defines the interface for LLM interactions.
type LLMClient interface {
	GeneratePlanSuggestion(ctx context.Context, req LLMPlanRequest) (*LLMPlanResponse, error)
	SuggestFix(ctx context.Context, req LLMFixRequest) (*LLMPlanResponse, error)
}

// LLMPlanRequest is the input for plan generation.
type LLMPlanRequest struct {
	Destination    string   `json:"destination"`
	NumDays        int      `json:"num_days"`
	Interests      []string `json:"interests"`
	Budget         string   `json:"budget"`
	TravelStyle    string   `json:"travel_style"`
	PlaceNames     []string `json:"place_names"`
	TotalBudgetYen *int     `json:"total_budget_yen,omitempty"`
}

// LLMFixRequest is the input for fixing violations.
type LLMFixRequest struct {
	CurrentPlan string   `json:"current_plan"`
	Violations  []string `json:"violations"`
}

// LLMPlanResponse is the structured response from the LLM.
type LLMPlanResponse struct {
	Days []LLMDayPlan `json:"days"`
}

// LLMDayPlan is a single day in the LLM response.
type LLMDayPlan struct {
	DayNumber  int           `json:"day_number"`
	Activities []LLMActivity `json:"activities"`
}

// LLMActivity is a single activity in the LLM response.
type LLMActivity struct {
	PlaceName        string `json:"place_name"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	DurationMin      int    `json:"duration_min"`
	Note             string `json:"note"`
	EstimatedCostYen int    `json:"estimated_cost_yen"`
}
