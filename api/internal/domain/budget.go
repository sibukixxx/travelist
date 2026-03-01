package domain

import "fmt"

// DayCost holds the total estimated cost for a single day.
type DayCost struct {
	DayNumber int `json:"day_number"`
	CostYen   int `json:"cost_yen"`
}

// BudgetSummary aggregates estimated costs across an itinerary.
type BudgetSummary struct {
	TotalCostYen int       `json:"total_cost_yen"`
	DailyCosts   []DayCost `json:"daily_costs"`
}

// CalcBudgetSummary computes a BudgetSummary from an Itinerary.
func CalcBudgetSummary(it *Itinerary) *BudgetSummary {
	dailyCosts := make([]DayCost, 0, len(it.Days))
	totalCost := 0

	for _, day := range it.Days {
		dayCost := 0
		for _, act := range day.Activities {
			dayCost += act.EstimatedCostYen
			if act.TravelFromPrev != nil {
				dayCost += act.TravelFromPrev.EstimatedCostYen
			}
		}
		dailyCosts = append(dailyCosts, DayCost{
			DayNumber: day.DayNumber,
			CostYen:   dayCost,
		})
		totalCost += dayCost
	}

	return &BudgetSummary{
		TotalCostYen: totalCost,
		DailyCosts:   dailyCosts,
	}
}

// ValidateBudget checks whether the total cost exceeds the given budget.
// If totalBudgetYen is nil, no validation is performed.
func ValidateBudget(summary *BudgetSummary, totalBudgetYen *int) []Violation {
	if totalBudgetYen == nil {
		return nil
	}
	if summary.TotalCostYen <= *totalBudgetYen {
		return nil
	}
	return []Violation{
		{
			Type:    ViolationBudgetExceeded,
			Message: fmt.Sprintf("total estimated cost %d円 exceeds budget %d円", summary.TotalCostYen, *totalBudgetYen),
		},
	}
}
