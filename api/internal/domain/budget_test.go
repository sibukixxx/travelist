package domain

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCalcBudgetSummary(t *testing.T) {
	baseDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	t.Run("returns zero totals for empty itinerary", func(t *testing.T) {
		it := &Itinerary{Days: []DayPlan{}}

		got := CalcBudgetSummary(it)

		want := &BudgetSummary{
			TotalCostYen: 0,
			DailyCosts:   []DayCost{},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("sums activity and travel costs for single day", func(t *testing.T) {
		it := &Itinerary{
			Days: []DayPlan{
				{
					DayNumber: 1,
					Date:      baseDate,
					Activities: []Activity{
						{
							Order:            0,
							StartTime:        "09:00",
							EndTime:          "10:30",
							EstimatedCostYen: 1500,
						},
						{
							Order:            1,
							StartTime:        "11:00",
							EndTime:          "12:00",
							EstimatedCostYen: 2000,
							TravelFromPrev: &TravelSegment{
								Mode:             TravelModeTrain,
								DurationMin:      15,
								DistanceM:        3000,
								EstimatedCostYen: 200,
							},
						},
					},
				},
			},
		}

		got := CalcBudgetSummary(it)

		want := &BudgetSummary{
			TotalCostYen: 3700,
			DailyCosts: []DayCost{
				{DayNumber: 1, CostYen: 3700},
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("sums costs across multiple days", func(t *testing.T) {
		it := &Itinerary{
			Days: []DayPlan{
				{
					DayNumber: 1,
					Date:      baseDate,
					Activities: []Activity{
						{EstimatedCostYen: 1000},
						{EstimatedCostYen: 500},
					},
				},
				{
					DayNumber: 2,
					Date:      baseDate.AddDate(0, 0, 1),
					Activities: []Activity{
						{
							EstimatedCostYen: 3000,
							TravelFromPrev: &TravelSegment{
								Mode:             TravelModeTaxi,
								EstimatedCostYen: 1500,
							},
						},
					},
				},
			},
		}

		got := CalcBudgetSummary(it)

		want := &BudgetSummary{
			TotalCostYen: 6000,
			DailyCosts: []DayCost{
				{DayNumber: 1, CostYen: 1500},
				{DayNumber: 2, CostYen: 4500},
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestValidateBudget(t *testing.T) {
	t.Run("returns no violation when budget is nil", func(t *testing.T) {
		summary := &BudgetSummary{TotalCostYen: 10000}

		got := ValidateBudget(summary, nil)

		if len(got) != 0 {
			t.Errorf("expected no violations, got %d", len(got))
		}
	})

	t.Run("returns no violation when within budget", func(t *testing.T) {
		summary := &BudgetSummary{TotalCostYen: 5000}
		budget := 10000

		got := ValidateBudget(summary, &budget)

		if len(got) != 0 {
			t.Errorf("expected no violations, got %d", len(got))
		}
	})

	t.Run("returns no violation when exactly at budget", func(t *testing.T) {
		summary := &BudgetSummary{TotalCostYen: 10000}
		budget := 10000

		got := ValidateBudget(summary, &budget)

		if len(got) != 0 {
			t.Errorf("expected no violations, got %d", len(got))
		}
	})

	t.Run("returns violation when over budget", func(t *testing.T) {
		summary := &BudgetSummary{TotalCostYen: 15000}
		budget := 10000

		got := ValidateBudget(summary, &budget)

		if len(got) != 1 {
			t.Fatalf("expected 1 violation, got %d", len(got))
		}
		if got[0].Type != ViolationBudgetExceeded {
			t.Errorf("expected type %s, got %s", ViolationBudgetExceeded, got[0].Type)
		}
	})
}
