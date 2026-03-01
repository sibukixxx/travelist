package clients_test

import (
	"context"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/infra/clients"
)

func TestStubLLMClient(t *testing.T) {
	t.Run("GeneratePlanSuggestion returns correct number of days", func(t *testing.T) {
		c := clients.NewStubLLMClient()
		ctx := context.Background()

		req := clients.LLMPlanRequest{
			Destination: "京都",
			NumDays:     3,
			Interests:   []string{"culture"},
			Budget:      "moderate",
			TravelStyle: "balanced",
			PlaceNames:  []string{"金閣寺", "伏見稲荷大社", "嵐山竹林", "清水寺", "錦市場"},
		}

		resp, err := c.GeneratePlanSuggestion(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Days) != 3 {
			t.Fatalf("expected 3 days, got %d", len(resp.Days))
		}
		for i, day := range resp.Days {
			if day.DayNumber != i+1 {
				t.Errorf("day %d: expected DayNumber %d, got %d", i, i+1, day.DayNumber)
			}
			if len(day.Activities) == 0 {
				t.Errorf("day %d: expected at least one activity", i+1)
			}
		}
	})

	t.Run("GeneratePlanSuggestion distributes places across days", func(t *testing.T) {
		c := clients.NewStubLLMClient()
		ctx := context.Background()

		req := clients.LLMPlanRequest{
			Destination: "京都",
			NumDays:     2,
			PlaceNames:  []string{"金閣寺", "伏見稲荷大社", "嵐山竹林", "清水寺"},
		}

		resp, err := c.GeneratePlanSuggestion(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		totalActivities := 0
		for _, day := range resp.Days {
			totalActivities += len(day.Activities)
		}
		if totalActivities != 4 {
			t.Errorf("expected 4 total activities, got %d", totalActivities)
		}
	})

	t.Run("GeneratePlanSuggestion activities have valid time fields", func(t *testing.T) {
		c := clients.NewStubLLMClient()
		ctx := context.Background()

		req := clients.LLMPlanRequest{
			Destination: "京都",
			NumDays:     1,
			PlaceNames:  []string{"金閣寺"},
		}

		resp, err := c.GeneratePlanSuggestion(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		act := resp.Days[0].Activities[0]
		if act.PlaceName != "金閣寺" {
			t.Errorf("expected PlaceName %q, got %q", "金閣寺", act.PlaceName)
		}
		if act.StartTime == "" || act.EndTime == "" {
			t.Error("StartTime or EndTime is empty")
		}
		if act.DurationMin <= 0 {
			t.Errorf("expected positive DurationMin, got %d", act.DurationMin)
		}
	})

	t.Run("SuggestFix returns a valid response", func(t *testing.T) {
		c := clients.NewStubLLMClient()
		ctx := context.Background()

		req := clients.LLMFixRequest{
			CurrentPlan: "some plan",
			Violations:  []string{"too_many_activities"},
		}

		resp, err := c.SuggestFix(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Days) == 0 {
			t.Fatal("expected at least one day in fix response")
		}
	})
}
