package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/clients"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

// --- Test Doubles (Stubs) ---

type stubPlaces struct {
	places        []domain.Place
	err           error
	capturedQuery string
}

func (s *stubPlaces) SearchPlaces(ctx context.Context, query string, lat, lng float64) ([]domain.Place, error) {
	s.capturedQuery = query
	return s.places, s.err
}

func (s *stubPlaces) GetPlaceDetails(ctx context.Context, placeID string) (*domain.Place, error) {
	return nil, nil
}

type stubLLM struct {
	resp        *clients.LLMPlanResponse
	err         error
	capturedReq *clients.LLMPlanRequest
}

func (s *stubLLM) GeneratePlanSuggestion(ctx context.Context, req clients.LLMPlanRequest) (*clients.LLMPlanResponse, error) {
	s.capturedReq = &req
	return s.resp, s.err
}

func (s *stubLLM) SuggestFix(ctx context.Context, req clients.LLMFixRequest) (*clients.LLMPlanResponse, error) {
	return nil, nil
}

type stubRepo struct {
	saved *domain.Itinerary
	err   error
}

func (s *stubRepo) Save(ctx context.Context, itinerary *domain.Itinerary) error {
	s.saved = itinerary
	return s.err
}

func (s *stubRepo) FindByID(ctx context.Context, id string) (*domain.Itinerary, error) {
	return nil, nil
}

func (s *stubRepo) List(ctx context.Context) ([]*domain.Itinerary, error) {
	return nil, nil
}

func (s *stubRepo) Delete(ctx context.Context, id string) error {
	return nil
}

// --- Helpers ---

func newFixedClock(t *testing.T) clock.FixedClock {
	t.Helper()
	return clock.FixedClock{Time: time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)}
}

func newDefaultRequest() domain.PlanRequest {
	return domain.PlanRequest{
		Destination: "京都",
		NumDays:     2,
		StartDate:   "2026-04-01",
		Preferences: domain.Preferences{
			Interests:   []string{"culture", "food"},
			Budget:      "moderate",
			TravelStyle: "balanced",
		},
		Constraint: domain.DefaultConstraint(),
	}
}

func newSamplePlaces() []domain.Place {
	return []domain.Place{
		{ID: "place-1", Name: "金閣寺", Lat: 35.0394, Lng: 135.7292},
		{ID: "place-2", Name: "伏見稲荷大社", Lat: 34.9671, Lng: 135.7727},
	}
}

func newSampleLLMResponse() *clients.LLMPlanResponse {
	return &clients.LLMPlanResponse{
		Days: []clients.LLMDayPlan{
			{
				DayNumber: 1,
				Activities: []clients.LLMActivity{
					{PlaceName: "金閣寺", StartTime: "09:00", EndTime: "11:00", DurationMin: 120, Note: "朝一で訪問"},
					{PlaceName: "伏見稲荷大社", StartTime: "13:00", EndTime: "15:00", DurationMin: 120, Note: "千本鳥居を散策"},
				},
			},
			{
				DayNumber: 2,
				Activities: []clients.LLMActivity{
					{PlaceName: "Unknown Spot", StartTime: "10:00", EndTime: "12:00", DurationMin: 120, Note: "隠れスポット"},
				},
			},
		},
	}
}

// --- Tests ---

func TestPlanGeneratorGenerate(t *testing.T) {
	t.Run("returns itinerary with correctly mapped places and dates", func(t *testing.T) {
		places := &stubPlaces{places: newSamplePlaces()}
		llm := &stubLLM{resp: newSampleLLMResponse()}
		repoStub := &stubRepo{}
		clk := newFixedClock(t)

		pg := usecase.NewPlanGenerator(places, llm, repoStub, clk)
		result, err := pg.Generate(context.Background(), newDefaultRequest())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		it := result.Itinerary

		// Title format
		if want := "京都 2日間の旅"; it.Title != want {
			t.Errorf("Title = %q, want %q", it.Title, want)
		}

		// Destination
		if it.Destination != "京都" {
			t.Errorf("Destination = %q, want %q", it.Destination, "京都")
		}

		// StartDate / EndDate
		wantStart := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
		wantEnd := time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC)
		if !it.StartDate.Equal(wantStart) {
			t.Errorf("StartDate = %v, want %v", it.StartDate, wantStart)
		}
		if !it.EndDate.Equal(wantEnd) {
			t.Errorf("EndDate = %v, want %v", it.EndDate, wantEnd)
		}

		// CreatedAt / UpdatedAt should use injected clock
		if !it.CreatedAt.Equal(clk.Time) {
			t.Errorf("CreatedAt = %v, want %v", it.CreatedAt, clk.Time)
		}
		if !it.UpdatedAt.Equal(clk.Time) {
			t.Errorf("UpdatedAt = %v, want %v", it.UpdatedAt, clk.Time)
		}

		// ID prefix
		if !strings.HasPrefix(it.ID, "itn_") {
			t.Errorf("ID = %q, want prefix \"itn_\"", it.ID)
		}

		// Days count
		if len(it.Days) != 2 {
			t.Fatalf("len(Days) = %d, want 2", len(it.Days))
		}

		// Day 1 metadata
		day1 := it.Days[0]
		if day1.DayNumber != 1 {
			t.Errorf("Days[0].DayNumber = %d, want 1", day1.DayNumber)
		}
		if !day1.Date.Equal(wantStart) {
			t.Errorf("Days[0].Date = %v, want %v", day1.Date, wantStart)
		}
		if len(day1.Activities) != 2 {
			t.Fatalf("len(Days[0].Activities) = %d, want 2", len(day1.Activities))
		}

		// Day 1, Activity 0 — known place "金閣寺" should be resolved
		act0 := day1.Activities[0]
		if act0.PlaceID != "place-1" {
			t.Errorf("act0.PlaceID = %q, want \"place-1\"", act0.PlaceID)
		}
		if act0.Place == nil || act0.Place.Name != "金閣寺" {
			t.Error("act0.Place should reference 金閣寺")
		}
		if act0.StartTime != "09:00" || act0.EndTime != "11:00" {
			t.Errorf("act0 time = %s-%s, want 09:00-11:00", act0.StartTime, act0.EndTime)
		}
		if act0.DurationMin != 120 {
			t.Errorf("act0.DurationMin = %d, want 120", act0.DurationMin)
		}
		if act0.Note != "朝一で訪問" {
			t.Errorf("act0.Note = %q, want %q", act0.Note, "朝一で訪問")
		}
		if act0.Order != 0 {
			t.Errorf("act0.Order = %d, want 0", act0.Order)
		}

		// Day 1, Activity 1 — known place "伏見稲荷大社"
		act1 := day1.Activities[1]
		if act1.PlaceID != "place-2" {
			t.Errorf("act1.PlaceID = %q, want \"place-2\"", act1.PlaceID)
		}
		if act1.Order != 1 {
			t.Errorf("act1.Order = %d, want 1", act1.Order)
		}

		// Day 2 date
		day2 := it.Days[1]
		if day2.DayNumber != 2 {
			t.Errorf("Days[1].DayNumber = %d, want 2", day2.DayNumber)
		}
		if !day2.Date.Equal(wantEnd) {
			t.Errorf("Days[1].Date = %v, want %v", day2.Date, wantEnd)
		}

		// Day 2, Activity 0 — unknown place should have empty PlaceID
		day2Act0 := day2.Activities[0]
		if day2Act0.PlaceID != "" {
			t.Errorf("day2Act0.PlaceID = %q, want empty (unknown place)", day2Act0.PlaceID)
		}
		if day2Act0.Place != nil {
			t.Errorf("day2Act0.Place = %v, want nil", day2Act0.Place)
		}

		// Itinerary was persisted
		if repoStub.saved == nil {
			t.Error("itinerary was not saved to repository")
		}
	})

	t.Run("passes destination and place names to LLM", func(t *testing.T) {
		places := &stubPlaces{places: newSamplePlaces()}
		llm := &stubLLM{resp: newSampleLLMResponse()}
		repoStub := &stubRepo{}

		pg := usecase.NewPlanGenerator(places, llm, repoStub, newFixedClock(t))
		_, err := pg.Generate(context.Background(), newDefaultRequest())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if llm.capturedReq == nil {
			t.Fatal("LLM was not called")
		}

		got := llm.capturedReq
		if got.Destination != "京都" {
			t.Errorf("LLM req.Destination = %q, want %q", got.Destination, "京都")
		}
		if got.NumDays != 2 {
			t.Errorf("LLM req.NumDays = %d, want 2", got.NumDays)
		}
		if got.Budget != "moderate" {
			t.Errorf("LLM req.Budget = %q, want %q", got.Budget, "moderate")
		}
		if got.TravelStyle != "balanced" {
			t.Errorf("LLM req.TravelStyle = %q, want %q", got.TravelStyle, "balanced")
		}

		wantPlaces := []string{"金閣寺", "伏見稲荷大社"}
		if diff := cmp.Diff(wantPlaces, got.PlaceNames); diff != "" {
			t.Errorf("LLM req.PlaceNames mismatch (-want +got):\n%s", diff)
		}
		wantInterests := []string{"culture", "food"}
		if diff := cmp.Diff(wantInterests, got.Interests); diff != "" {
			t.Errorf("LLM req.Interests mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("searches places with destination query", func(t *testing.T) {
		places := &stubPlaces{places: newSamplePlaces()}
		llm := &stubLLM{resp: newSampleLLMResponse()}
		repoStub := &stubRepo{}

		pg := usecase.NewPlanGenerator(places, llm, repoStub, newFixedClock(t))
		_, err := pg.Generate(context.Background(), newDefaultRequest())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if places.capturedQuery != "京都" {
			t.Errorf("places query = %q, want %q", places.capturedQuery, "京都")
		}
	})

	t.Run("returns violations without failing when constraints are violated", func(t *testing.T) {
		llmResp := &clients.LLMPlanResponse{
			Days: []clients.LLMDayPlan{
				{
					DayNumber: 1,
					Activities: []clients.LLMActivity{
						{PlaceName: "早朝スポット", StartTime: "06:00", EndTime: "07:00", DurationMin: 60},
					},
				},
			},
		}

		places := &stubPlaces{places: []domain.Place{}}
		llm := &stubLLM{resp: llmResp}
		repoStub := &stubRepo{}

		pg := usecase.NewPlanGenerator(places, llm, repoStub, newFixedClock(t))
		req := newDefaultRequest()
		req.NumDays = 1
		result, err := pg.Generate(context.Background(), req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Itinerary == nil {
			t.Fatal("expected itinerary, got nil")
		}
		if repoStub.saved == nil {
			t.Error("itinerary should be saved even with violations")
		}

		hasOutsideHours := false
		for _, v := range result.Violations {
			if v.Type == domain.ViolationOutsideHours {
				hasOutsideHours = true
			}
		}
		if !hasOutsideHours {
			t.Errorf("expected ViolationOutsideHours, got %v", result.Violations)
		}
	})

	t.Run("returns no violations when plan is valid", func(t *testing.T) {
		places := &stubPlaces{places: newSamplePlaces()}
		llm := &stubLLM{resp: newSampleLLMResponse()}
		repoStub := &stubRepo{}

		pg := usecase.NewPlanGenerator(places, llm, repoStub, newFixedClock(t))
		result, err := pg.Generate(context.Background(), newDefaultRequest())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Violations) != 0 {
			t.Errorf("expected no violations, got %v", result.Violations)
		}
	})
}

func TestPlanGeneratorGenerateErrors(t *testing.T) {
	tests := []struct {
		name    string
		req     domain.PlanRequest
		places  *stubPlaces
		llm     *stubLLM
		repo    *stubRepo
		wantErr string
	}{
		{
			name: "returns error when start_date is invalid",
			req: func() domain.PlanRequest {
				r := newDefaultRequest()
				r.StartDate = "not-a-date"
				return r
			}(),
			places:  &stubPlaces{places: newSamplePlaces()},
			llm:     &stubLLM{resp: newSampleLLMResponse()},
			repo:    &stubRepo{},
			wantErr: "invalid start_date",
		},
		{
			name:    "returns error when places search fails",
			req:     newDefaultRequest(),
			places:  &stubPlaces{err: errors.New("API unavailable")},
			llm:     &stubLLM{resp: newSampleLLMResponse()},
			repo:    &stubRepo{},
			wantErr: "places search failed",
		},
		{
			name:    "returns error when LLM generation fails",
			req:     newDefaultRequest(),
			places:  &stubPlaces{places: newSamplePlaces()},
			llm:     &stubLLM{err: errors.New("rate limit exceeded")},
			repo:    &stubRepo{},
			wantErr: "LLM plan generation failed",
		},
		{
			name:    "returns error when repository save fails",
			req:     newDefaultRequest(),
			places:  &stubPlaces{places: newSamplePlaces()},
			llm:     &stubLLM{resp: newSampleLLMResponse()},
			repo:    &stubRepo{err: errors.New("connection refused")},
			wantErr: "save failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := usecase.NewPlanGenerator(tt.places, tt.llm, tt.repo, newFixedClock(t))
			_, err := pg.Generate(context.Background(), tt.req)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}
