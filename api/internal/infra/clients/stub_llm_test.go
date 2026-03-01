package clients_test

import (
	"context"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/infra/clients"
)

func TestStubLLMClientGeneratePlanSuggestion(t *testing.T) {
	t.Run("returns days matching requested num_days", func(t *testing.T) {
		client := clients.NewStubLLMClient()
		resp, err := client.GeneratePlanSuggestion(context.Background(), clients.LLMPlanRequest{
			Destination: "京都",
			NumDays:     3,
			Interests:   []string{"culture"},
			Budget:      "moderate",
			TravelStyle: "balanced",
			PlaceNames:  []string{"金閣寺", "伏見稲荷大社", "清水寺"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Days) != 3 {
			t.Fatalf("expected 3 days, got %d", len(resp.Days))
		}
		for i, day := range resp.Days {
			if day.DayNumber != i+1 {
				t.Errorf("Days[%d].DayNumber = %d, want %d", i, day.DayNumber, i+1)
			}
			if len(day.Activities) == 0 {
				t.Errorf("Days[%d] has no activities", i)
			}
		}
	})

	t.Run("activities have valid time format", func(t *testing.T) {
		client := clients.NewStubLLMClient()
		resp, err := client.GeneratePlanSuggestion(context.Background(), clients.LLMPlanRequest{
			Destination: "東京",
			NumDays:     1,
			PlaceNames:  []string{"東京タワー"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, day := range resp.Days {
			for _, act := range day.Activities {
				if len(act.StartTime) != 5 || act.StartTime[2] != ':' {
					t.Errorf("invalid StartTime format: %q", act.StartTime)
				}
				if len(act.EndTime) != 5 || act.EndTime[2] != ':' {
					t.Errorf("invalid EndTime format: %q", act.EndTime)
				}
				if act.DurationMin <= 0 {
					t.Errorf("DurationMin should be positive, got %d", act.DurationMin)
				}
			}
		}
	})

	t.Run("uses place names from request when available", func(t *testing.T) {
		client := clients.NewStubLLMClient()
		placeNames := []string{"金閣寺", "伏見稲荷大社"}
		resp, err := client.GeneratePlanSuggestion(context.Background(), clients.LLMPlanRequest{
			Destination: "京都",
			NumDays:     1,
			PlaceNames:  placeNames,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// At least one activity should reference a requested place name
		found := false
		for _, day := range resp.Days {
			for _, act := range day.Activities {
				for _, name := range placeNames {
					if act.PlaceName == name {
						found = true
					}
				}
			}
		}
		if !found {
			t.Error("expected at least one activity to use a place name from request")
		}
	})
}

func TestStubLLMClientSuggestFix(t *testing.T) {
	t.Run("returns non-empty response", func(t *testing.T) {
		client := clients.NewStubLLMClient()
		resp, err := client.SuggestFix(context.Background(), clients.LLMFixRequest{
			CurrentPlan: "some plan",
			Violations:  []string{"time overlap"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Days) == 0 {
			t.Error("expected non-empty days in fix response")
		}
	})
}

// Verify interface compliance at compile time.
var _ clients.LLMClient = (*clients.StubLLMClient)(nil)
