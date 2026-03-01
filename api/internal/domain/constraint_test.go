package domain_test

import (
	"testing"
	"time"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

func TestValidateDayPlan_NoViolations(t *testing.T) {
	day := domain.DayPlan{
		DayNumber: 1,
		Date:      time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), // Wednesday
		Activities: []domain.Activity{
			{Order: 0, StartTime: "09:00", EndTime: "11:00", DurationMin: 120},
			{Order: 1, StartTime: "11:30", EndTime: "13:00", DurationMin: 90},
		},
	}
	violations := domain.ValidateDayPlan(day, domain.DefaultConstraint())
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestValidateDayPlan_TooManyActivities(t *testing.T) {
	activities := make([]domain.Activity, 7)
	for i := range activities {
		h := 8 + i*2
		activities[i] = domain.Activity{
			Order:       i,
			StartTime:   formatTime(h, 0),
			EndTime:     formatTime(h+1, 0),
			DurationMin: 60,
		}
	}
	day := domain.DayPlan{
		DayNumber:  1,
		Date:       time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		Activities: activities,
	}
	violations := domain.ValidateDayPlan(day, domain.DefaultConstraint())
	hasViolation := false
	for _, v := range violations {
		if v.Type == domain.ViolationTooManyActivities {
			hasViolation = true
		}
	}
	if !hasViolation {
		t.Error("expected TooManyActivities violation")
	}
}

func TestValidateDayPlan_OutsideHours(t *testing.T) {
	day := domain.DayPlan{
		DayNumber: 1,
		Date:      time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		Activities: []domain.Activity{
			{Order: 0, StartTime: "06:00", EndTime: "08:00", DurationMin: 120},
		},
	}
	violations := domain.ValidateDayPlan(day, domain.DefaultConstraint())
	hasViolation := false
	for _, v := range violations {
		if v.Type == domain.ViolationOutsideHours {
			hasViolation = true
		}
	}
	if !hasViolation {
		t.Error("expected OutsideHours violation")
	}
}

func TestValidateDayPlan_TimeOverlap(t *testing.T) {
	day := domain.DayPlan{
		DayNumber: 1,
		Date:      time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		Activities: []domain.Activity{
			{Order: 0, StartTime: "09:00", EndTime: "11:00", DurationMin: 120},
			{Order: 1, StartTime: "10:30", EndTime: "12:00", DurationMin: 90},
		},
	}
	violations := domain.ValidateDayPlan(day, domain.DefaultConstraint())
	hasViolation := false
	for _, v := range violations {
		if v.Type == domain.ViolationTimeOverlap {
			hasViolation = true
		}
	}
	if !hasViolation {
		t.Error("expected TimeOverlap violation")
	}
}

func TestValidateDayPlan_ExcessiveWalk(t *testing.T) {
	day := domain.DayPlan{
		DayNumber: 1,
		Date:      time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		Activities: []domain.Activity{
			{Order: 0, StartTime: "09:00", EndTime: "11:00", DurationMin: 120},
			{
				Order: 1, StartTime: "11:30", EndTime: "13:00", DurationMin: 90,
				TravelFromPrev: &domain.TravelSegment{
					Mode:        domain.TravelModeWalk,
					DurationMin: 45,
					DistanceM:   3500,
				},
			},
		},
	}
	violations := domain.ValidateDayPlan(day, domain.DefaultConstraint())
	hasViolation := false
	for _, v := range violations {
		if v.Type == domain.ViolationExcessiveWalk {
			hasViolation = true
		}
	}
	if !hasViolation {
		t.Error("expected ExcessiveWalk violation")
	}
}

func TestValidateDayPlan_ClosedPlace(t *testing.T) {
	place := &domain.Place{
		Name: "Museum",
		OpeningHours: &domain.OpeningHours{
			Periods: []domain.Period{
				{DayOfWeek: 3, OpenTime: "10:00", CloseTime: "17:00"}, // Wednesday only
			},
		},
	}
	day := domain.DayPlan{
		DayNumber: 1,
		Date:      time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), // Wednesday
		Activities: []domain.Activity{
			{Order: 0, StartTime: "08:30", EndTime: "10:00", DurationMin: 90, Place: place},
		},
	}
	violations := domain.ValidateDayPlan(day, domain.DefaultConstraint())
	hasViolation := false
	for _, v := range violations {
		if v.Type == domain.ViolationClosedPlace {
			hasViolation = true
		}
	}
	if !hasViolation {
		t.Error("expected ClosedPlace violation")
	}
}

func formatTime(h, m int) string {
	return time.Date(0, 1, 1, h, m, 0, 0, time.UTC).Format("15:04")
}
