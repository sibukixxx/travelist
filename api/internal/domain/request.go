package domain

// PlanRequest represents a user's request to generate a travel plan.
type PlanRequest struct {
	Destination string      `json:"destination"`
	NumDays     int         `json:"num_days"`
	StartDate   string      `json:"start_date"` // "YYYY-MM-DD"
	Preferences Preferences `json:"preferences"`
	Constraint  Constraint  `json:"constraint"`
}

// Preferences represents user preferences for the trip.
type Preferences struct {
	Interests      []string `json:"interests"`        // e.g. ["culture", "food", "nature"]
	Budget         string   `json:"budget"`           // "budget", "moderate", "luxury"
	TravelStyle    string   `json:"travel_style"`     // "relaxed", "active", "balanced"
	TotalBudgetYen *int     `json:"total_budget_yen"` // optional numeric budget cap in JPY
}
