package clients

import (
	"context"
	"fmt"
)

// StubLLMClient is an in-memory stub implementation of LLMClient.
type StubLLMClient struct{}

// NewStubLLMClient creates a new StubLLMClient.
func NewStubLLMClient() *StubLLMClient {
	return &StubLLMClient{}
}

func (s *StubLLMClient) GeneratePlanSuggestion(_ context.Context, req LLMPlanRequest) (*LLMPlanResponse, error) {
	days := make([]LLMDayPlan, req.NumDays)

	// Schedule templates: morning, lunch, afternoon activities
	templates := []struct {
		startTime   string
		endTime     string
		durationMin int
		notePrefix  string
		costYen     int
	}{
		{"09:00", "11:00", 120, "午前の観光", 500},
		{"11:30", "12:30", 60, "ランチ", 1500},
		{"13:30", "15:30", 120, "午後の観光", 800},
	}

	for d := 0; d < req.NumDays; d++ {
		day := LLMDayPlan{DayNumber: d + 1}

		for i, tmpl := range templates {
			// Pick place name: use from request if available, otherwise generate
			placeName := fmt.Sprintf("%s 観光スポット%d", req.Destination, d*3+i+1)
			placeIdx := d*len(templates) + i
			if placeIdx < len(req.PlaceNames) {
				placeName = req.PlaceNames[placeIdx]
			}

			day.Activities = append(day.Activities, LLMActivity{
				PlaceName:        placeName,
				StartTime:        tmpl.startTime,
				EndTime:          tmpl.endTime,
				DurationMin:      tmpl.durationMin,
				Note:             fmt.Sprintf("%s: %s", tmpl.notePrefix, placeName),
				EstimatedCostYen: tmpl.costYen,
			})
		}

		days[d] = day
	}

	return &LLMPlanResponse{Days: days}, nil
}

func (s *StubLLMClient) SuggestFix(_ context.Context, req LLMFixRequest) (*LLMPlanResponse, error) {
	return &LLMPlanResponse{
		Days: []LLMDayPlan{
			{
				DayNumber: 1,
				Activities: []LLMActivity{
					{
						PlaceName:        "修正済みスポット",
						StartTime:        "10:00",
						EndTime:          "12:00",
						DurationMin:      120,
						Note:             "修正プラン",
						EstimatedCostYen: 1000,
					},
				},
			},
		},
	}, nil
}
