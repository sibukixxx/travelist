package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/clients"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
)

// PlanGenerator handles the itinerary generation workflow.
type PlanGenerator struct {
	places PlacesClientInterface
	llm    LLMClientInterface
	repo   repo.ItineraryRepository
	clock  clock.Clock
}

// PlacesClientInterface is used by the usecase layer.
type PlacesClientInterface = clients.PlacesClient

// LLMClientInterface is used by the usecase layer.
type LLMClientInterface = clients.LLMClient

// NewPlanGenerator creates a new PlanGenerator.
func NewPlanGenerator(
	places PlacesClientInterface,
	llm LLMClientInterface,
	repository repo.ItineraryRepository,
	clk clock.Clock,
) *PlanGenerator {
	return &PlanGenerator{
		places: places,
		llm:    llm,
		repo:   repository,
		clock:  clk,
	}
}

// GenerateResult holds the result of plan generation.
type GenerateResult struct {
	Itinerary  *domain.Itinerary  `json:"itinerary"`
	Violations []domain.Violation `json:"violations"`
}

// Generate creates a new itinerary based on the request.
func (pg *PlanGenerator) Generate(ctx context.Context, req domain.PlanRequest) (*GenerateResult, error) {
	// 1. Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	// 2. Search for candidate places
	candidatePlaces, err := pg.places.SearchPlaces(ctx, req.Destination, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("places search failed: %w", err)
	}

	// 3. Build place name list for LLM
	placeNames := make([]string, len(candidatePlaces))
	for i, p := range candidatePlaces {
		placeNames[i] = p.Name
	}

	// 4. Ask LLM for a draft plan
	llmResp, err := pg.llm.GeneratePlanSuggestion(ctx, clients.LLMPlanRequest{
		Destination: req.Destination,
		NumDays:     req.NumDays,
		Interests:   req.Preferences.Interests,
		Budget:      req.Preferences.Budget,
		TravelStyle: req.Preferences.TravelStyle,
		PlaceNames:  placeNames,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM plan generation failed: %w", err)
	}

	// 5. Convert LLM response to domain model
	itinerary := &domain.Itinerary{
		ID:          generateID(),
		Title:       fmt.Sprintf("%s %d日間の旅", req.Destination, req.NumDays),
		Destination: req.Destination,
		StartDate:   startDate,
		EndDate:     startDate.AddDate(0, 0, req.NumDays-1),
		CreatedAt:   pg.clock.Now(),
		UpdatedAt:   pg.clock.Now(),
	}

	placeByName := make(map[string]domain.Place)
	for _, p := range candidatePlaces {
		placeByName[p.Name] = p
	}

	for _, llmDay := range llmResp.Days {
		dayPlan := domain.DayPlan{
			DayNumber: llmDay.DayNumber,
			Date:      startDate.AddDate(0, 0, llmDay.DayNumber-1),
		}
		for j, llmAct := range llmDay.Activities {
			act := domain.Activity{
				Order:       j,
				StartTime:   llmAct.StartTime,
				EndTime:     llmAct.EndTime,
				DurationMin: llmAct.DurationMin,
				Note:        llmAct.Note,
			}
			if p, ok := placeByName[llmAct.PlaceName]; ok {
				act.PlaceID = p.ID
				act.Place = &p
			}
			dayPlan.Activities = append(dayPlan.Activities, act)
		}
		itinerary.Days = append(itinerary.Days, dayPlan)
	}

	// 6. Validate with rule-based logic
	var allViolations []domain.Violation
	for _, day := range itinerary.Days {
		violations := domain.ValidateDayPlan(day, req.Constraint)
		allViolations = append(allViolations, violations...)
	}

	// 7. Save
	if err := pg.repo.Save(ctx, itinerary); err != nil {
		return nil, fmt.Errorf("save failed: %w", err)
	}

	return &GenerateResult{
		Itinerary:  itinerary,
		Violations: allViolations,
	}, nil
}

func generateID() string {
	return fmt.Sprintf("itn_%d", time.Now().UnixNano())
}
