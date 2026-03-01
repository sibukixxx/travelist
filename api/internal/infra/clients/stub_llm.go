package clients

import (
	"context"
	"fmt"
)

// StubLLMClient is a development stub that generates deterministic plan suggestions.
type StubLLMClient struct{}

// NewStubLLMClient returns a new StubLLMClient.
func NewStubLLMClient() *StubLLMClient {
	return &StubLLMClient{}
}

func (c *StubLLMClient) GeneratePlanSuggestion(_ context.Context, req LLMPlanRequest) (*LLMPlanResponse, error) {
	days := make([]LLMDayPlan, req.NumDays)

	// Distribute places round-robin across days
	for i := 0; i < req.NumDays; i++ {
		days[i] = LLMDayPlan{
			DayNumber:  i + 1,
			Activities: []LLMActivity{},
		}
	}

	for i, name := range req.PlaceNames {
		dayIdx := i % req.NumDays
		order := len(days[dayIdx].Activities)
		startHour := 9 + order*2 // 09:00, 11:00, 13:00, ...
		days[dayIdx].Activities = append(days[dayIdx].Activities, LLMActivity{
			PlaceName:        name,
			StartTime:        fmt.Sprintf("%02d:00", startHour),
			EndTime:          fmt.Sprintf("%02d:30", startHour+1),
			DurationMin:      90,
			Note:             fmt.Sprintf("%sを観光", name),
			EstimatedCostYen: 500,
		})
	}

	// Ensure every day has at least one activity
	for i := range days {
		if len(days[i].Activities) == 0 {
			days[i].Activities = append(days[i].Activities, LLMActivity{
				PlaceName:        fmt.Sprintf("%s散策", req.Destination),
				StartTime:        "09:00",
				EndTime:          "10:30",
				DurationMin:      90,
				Note:             "周辺を散策",
				EstimatedCostYen: 0,
			})
		}
	}

	return &LLMPlanResponse{Days: days}, nil
}

func (c *StubLLMClient) SuggestFix(_ context.Context, _ LLMFixRequest) (*LLMPlanResponse, error) {
	return &LLMPlanResponse{
		Days: []LLMDayPlan{
			{
				DayNumber: 1,
				Activities: []LLMActivity{
					{
						PlaceName:        "修正済みプラン",
						StartTime:        "09:00",
						EndTime:          "10:30",
						DurationMin:      90,
						Note:             "制約違反を修正した代替プラン",
						EstimatedCostYen: 500,
					},
				},
			},
		},
	}, nil
}
