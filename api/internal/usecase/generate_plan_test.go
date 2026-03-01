package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/clients"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
)

// --- Fakes ---

type fakePlaces struct {
	places []domain.Place
	err    error
}

func (f *fakePlaces) SearchPlaces(_ context.Context, _ string, _, _ float64) ([]domain.Place, error) {
	return f.places, f.err
}

func (f *fakePlaces) GetPlaceDetails(_ context.Context, _ string) (*domain.Place, error) {
	return nil, nil
}

type fakeLLM struct {
	resp *clients.LLMPlanResponse
	err  error
}

func (f *fakeLLM) GeneratePlanSuggestion(_ context.Context, _ clients.LLMPlanRequest) (*clients.LLMPlanResponse, error) {
	return f.resp, f.err
}

func (f *fakeLLM) SuggestFix(_ context.Context, _ clients.LLMFixRequest) (*clients.LLMPlanResponse, error) {
	return nil, nil
}

type fakeRepo struct {
	saved *domain.Itinerary
	err   error
}

func (f *fakeRepo) Save(_ context.Context, it *domain.Itinerary) error {
	f.saved = it
	return f.err
}

func (f *fakeRepo) FindByID(_ context.Context, _ string) (*domain.Itinerary, error) {
	return nil, nil
}

func (f *fakeRepo) List(_ context.Context) ([]*domain.Itinerary, error) {
	return nil, nil
}

func (f *fakeRepo) Delete(_ context.Context, _ string) error {
	return nil
}

// --- Helpers ---

var fixedTime = time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)

func newTestGenerator(places *fakePlaces, llm *fakeLLM, repository *fakeRepo) *PlanGenerator {
	return NewPlanGenerator(places, llm, repository, clock.FixedClock{Time: fixedTime})
}

func validRequest() domain.PlanRequest {
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

func samplePlaces() []domain.Place {
	return []domain.Place{
		{ID: "p1", Name: "金閣寺", Rating: 4.5},
		{ID: "p2", Name: "伏見稲荷大社", Rating: 4.7},
		{ID: "p3", Name: "錦市場", Rating: 4.3},
	}
}

func sampleLLMResponse() *clients.LLMPlanResponse {
	return &clients.LLMPlanResponse{
		Days: []clients.LLMDayPlan{
			{
				DayNumber: 1,
				Activities: []clients.LLMActivity{
					{PlaceName: "金閣寺", StartTime: "09:00", EndTime: "11:00", DurationMin: 120, Note: "朝一番で訪問"},
					{PlaceName: "錦市場", StartTime: "12:00", EndTime: "14:00", DurationMin: 120, Note: "昼食"},
				},
			},
			{
				DayNumber: 2,
				Activities: []clients.LLMActivity{
					{PlaceName: "伏見稲荷大社", StartTime: "09:00", EndTime: "12:00", DurationMin: 180, Note: "千本鳥居"},
				},
			},
		},
	}
}

// --- Tests ---

func TestGenerate_Success(t *testing.T) {
	// Arrange
	places := &fakePlaces{places: samplePlaces()}
	llm := &fakeLLM{resp: sampleLLMResponse()}
	repository := &fakeRepo{}
	pg := newTestGenerator(places, llm, repository)

	// Act
	result, err := pg.Generate(context.Background(), validRequest())

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Itinerary == nil {
		t.Fatal("expected itinerary, got nil")
	}
	it := result.Itinerary
	if it.Destination != "京都" {
		t.Errorf("destination = %q, want %q", it.Destination, "京都")
	}
	if len(it.Days) != 2 {
		t.Fatalf("days = %d, want 2", len(it.Days))
	}
	if len(it.Days[0].Activities) != 2 {
		t.Errorf("day1 activities = %d, want 2", len(it.Days[0].Activities))
	}
	if len(it.Days[1].Activities) != 1 {
		t.Errorf("day2 activities = %d, want 1", len(it.Days[1].Activities))
	}
	// Verify itinerary was saved
	if repository.saved == nil {
		t.Error("expected itinerary to be saved")
	}
}

func TestGenerate_TitleFormat(t *testing.T) {
	// Arrange
	places := &fakePlaces{places: samplePlaces()}
	llm := &fakeLLM{resp: sampleLLMResponse()}
	pg := newTestGenerator(places, llm, &fakeRepo{})

	// Act
	result, err := pg.Generate(context.Background(), validRequest())

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "京都 2日間の旅"
	if result.Itinerary.Title != want {
		t.Errorf("title = %q, want %q", result.Itinerary.Title, want)
	}
}

func TestGenerate_DateCalculation(t *testing.T) {
	// Arrange
	places := &fakePlaces{places: samplePlaces()}
	llm := &fakeLLM{resp: sampleLLMResponse()}
	pg := newTestGenerator(places, llm, &fakeRepo{})

	// Act
	result, err := pg.Generate(context.Background(), validRequest())

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	it := result.Itinerary
	wantStart := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC)
	if !it.StartDate.Equal(wantStart) {
		t.Errorf("start_date = %v, want %v", it.StartDate, wantStart)
	}
	if !it.EndDate.Equal(wantEnd) {
		t.Errorf("end_date = %v, want %v", it.EndDate, wantEnd)
	}
	// Day dates
	if !it.Days[0].Date.Equal(wantStart) {
		t.Errorf("day1 date = %v, want %v", it.Days[0].Date, wantStart)
	}
	if !it.Days[1].Date.Equal(wantEnd) {
		t.Errorf("day2 date = %v, want %v", it.Days[1].Date, wantEnd)
	}
}

func TestGenerate_TimestampsUseFixedClock(t *testing.T) {
	// Arrange
	places := &fakePlaces{places: samplePlaces()}
	llm := &fakeLLM{resp: sampleLLMResponse()}
	pg := newTestGenerator(places, llm, &fakeRepo{})

	// Act
	result, err := pg.Generate(context.Background(), validRequest())

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Itinerary.CreatedAt.Equal(fixedTime) {
		t.Errorf("created_at = %v, want %v", result.Itinerary.CreatedAt, fixedTime)
	}
	if !result.Itinerary.UpdatedAt.Equal(fixedTime) {
		t.Errorf("updated_at = %v, want %v", result.Itinerary.UpdatedAt, fixedTime)
	}
}

func TestGenerate_PlaceMatching(t *testing.T) {
	// Arrange
	places := &fakePlaces{places: samplePlaces()}
	llm := &fakeLLM{resp: sampleLLMResponse()}
	pg := newTestGenerator(places, llm, &fakeRepo{})

	// Act
	result, err := pg.Generate(context.Background(), validRequest())

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	act0 := result.Itinerary.Days[0].Activities[0]
	if act0.PlaceID != "p1" {
		t.Errorf("activity 0 place_id = %q, want %q", act0.PlaceID, "p1")
	}
	if act0.Place == nil || act0.Place.Name != "金閣寺" {
		t.Errorf("activity 0 place name = %v, want 金閣寺", act0.Place)
	}
}

func TestGenerate_UnknownPlaceNotLinked(t *testing.T) {
	// Arrange: LLM returns a place name that doesn't exist in candidates
	places := &fakePlaces{places: samplePlaces()}
	llm := &fakeLLM{resp: &clients.LLMPlanResponse{
		Days: []clients.LLMDayPlan{
			{
				DayNumber: 1,
				Activities: []clients.LLMActivity{
					{PlaceName: "存在しない場所", StartTime: "09:00", EndTime: "11:00", DurationMin: 120},
				},
			},
		},
	}}
	req := validRequest()
	req.NumDays = 1
	pg := newTestGenerator(places, llm, &fakeRepo{})

	// Act
	result, err := pg.Generate(context.Background(), req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	act := result.Itinerary.Days[0].Activities[0]
	if act.PlaceID != "" {
		t.Errorf("expected empty place_id for unknown place, got %q", act.PlaceID)
	}
	if act.Place != nil {
		t.Errorf("expected nil place for unknown place, got %v", act.Place)
	}
}

func TestGenerate_InvalidStartDate(t *testing.T) {
	// Arrange
	pg := newTestGenerator(&fakePlaces{places: samplePlaces()}, &fakeLLM{resp: sampleLLMResponse()}, &fakeRepo{})
	req := validRequest()
	req.StartDate = "not-a-date"

	// Act
	_, err := pg.Generate(context.Background(), req)

	// Assert
	if err == nil {
		t.Fatal("expected error for invalid start_date")
	}
	if !strings.Contains(err.Error(), "invalid start_date") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "invalid start_date")
	}
}

func TestGenerate_PlacesSearchError(t *testing.T) {
	// Arrange
	placesErr := errors.New("API quota exceeded")
	pg := newTestGenerator(&fakePlaces{err: placesErr}, &fakeLLM{resp: sampleLLMResponse()}, &fakeRepo{})

	// Act
	_, err := pg.Generate(context.Background(), validRequest())

	// Assert
	if err == nil {
		t.Fatal("expected error when places search fails")
	}
	if !errors.Is(err, placesErr) {
		t.Errorf("error should wrap original: %v", err)
	}
	if !strings.Contains(err.Error(), "places search failed") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "places search failed")
	}
}

func TestGenerate_LLMError(t *testing.T) {
	// Arrange
	llmErr := errors.New("model overloaded")
	pg := newTestGenerator(&fakePlaces{places: samplePlaces()}, &fakeLLM{err: llmErr}, &fakeRepo{})

	// Act
	_, err := pg.Generate(context.Background(), validRequest())

	// Assert
	if err == nil {
		t.Fatal("expected error when LLM fails")
	}
	if !errors.Is(err, llmErr) {
		t.Errorf("error should wrap original: %v", err)
	}
	if !strings.Contains(err.Error(), "LLM plan generation failed") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "LLM plan generation failed")
	}
}

func TestGenerate_SaveError(t *testing.T) {
	// Arrange
	saveErr := errors.New("database connection lost")
	places := &fakePlaces{places: samplePlaces()}
	llm := &fakeLLM{resp: sampleLLMResponse()}
	repository := &fakeRepo{err: saveErr}
	pg := newTestGenerator(places, llm, repository)

	// Act
	_, err := pg.Generate(context.Background(), validRequest())

	// Assert
	if err == nil {
		t.Fatal("expected error when save fails")
	}
	if !errors.Is(err, saveErr) {
		t.Errorf("error should wrap original: %v", err)
	}
	if !strings.Contains(err.Error(), "save failed") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "save failed")
	}
}

func TestGenerate_ViolationsReturnedButSaveSucceeds(t *testing.T) {
	// Arrange: activity starts before earliest allowed time → OutsideHours violation
	places := &fakePlaces{places: samplePlaces()}
	llm := &fakeLLM{resp: &clients.LLMPlanResponse{
		Days: []clients.LLMDayPlan{
			{
				DayNumber: 1,
				Activities: []clients.LLMActivity{
					{PlaceName: "金閣寺", StartTime: "06:00", EndTime: "08:00", DurationMin: 120},
				},
			},
		},
	}}
	req := validRequest()
	req.NumDays = 1
	repository := &fakeRepo{}
	pg := newTestGenerator(places, llm, repository)

	// Act
	result, err := pg.Generate(context.Background(), req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Violations) == 0 {
		t.Error("expected violations but got none")
	}
	hasOutsideHours := false
	for _, v := range result.Violations {
		if v.Type == domain.ViolationOutsideHours {
			hasOutsideHours = true
		}
	}
	if !hasOutsideHours {
		t.Error("expected OutsideHours violation")
	}
	// Verify itinerary was still saved despite violations
	if repository.saved == nil {
		t.Error("expected itinerary to be saved even with violations")
	}
}

func TestGenerate_ActivityFieldsMapping(t *testing.T) {
	// Arrange
	places := &fakePlaces{places: samplePlaces()}
	llm := &fakeLLM{resp: sampleLLMResponse()}
	pg := newTestGenerator(places, llm, &fakeRepo{})

	// Act
	result, err := pg.Generate(context.Background(), validRequest())

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	act := result.Itinerary.Days[0].Activities[0]
	if act.StartTime != "09:00" {
		t.Errorf("start_time = %q, want %q", act.StartTime, "09:00")
	}
	if act.EndTime != "11:00" {
		t.Errorf("end_time = %q, want %q", act.EndTime, "11:00")
	}
	if act.DurationMin != 120 {
		t.Errorf("duration_min = %d, want %d", act.DurationMin, 120)
	}
	if act.Note != "朝一番で訪問" {
		t.Errorf("note = %q, want %q", act.Note, "朝一番で訪問")
	}
	if act.Order != 0 {
		t.Errorf("order = %d, want %d", act.Order, 0)
	}
}

func TestGenerate_EmptyPlacesStillCallsLLM(t *testing.T) {
	// Arrange: places returns empty list
	places := &fakePlaces{places: []domain.Place{}}
	llm := &fakeLLM{resp: &clients.LLMPlanResponse{
		Days: []clients.LLMDayPlan{
			{
				DayNumber: 1,
				Activities: []clients.LLMActivity{
					{PlaceName: "Some Place", StartTime: "10:00", EndTime: "12:00", DurationMin: 120},
				},
			},
		},
	}}
	req := validRequest()
	req.NumDays = 1
	pg := newTestGenerator(places, llm, &fakeRepo{})

	// Act
	result, err := pg.Generate(context.Background(), req)

	// Assert: should still succeed; activities just won't have place links
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Itinerary == nil {
		t.Fatal("expected itinerary, got nil")
	}
	act := result.Itinerary.Days[0].Activities[0]
	if act.PlaceID != "" {
		t.Errorf("expected empty place_id, got %q", act.PlaceID)
	}
}
