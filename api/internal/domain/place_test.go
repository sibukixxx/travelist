package domain_test

import (
	"testing"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

func TestPlace_IsOpenAt_NoHours(t *testing.T) {
	p := &domain.Place{Name: "Park"}
	if !p.IsOpenAt(0, "12:00") {
		t.Error("place with no opening hours should be assumed open")
	}
}

func TestPlace_IsOpenAt_Open(t *testing.T) {
	p := &domain.Place{
		Name: "Cafe",
		OpeningHours: &domain.OpeningHours{
			Periods: []domain.Period{
				{DayOfWeek: 1, OpenTime: "08:00", CloseTime: "20:00"},
			},
		},
	}
	if !p.IsOpenAt(1, "12:00") {
		t.Error("expected place to be open on Monday at 12:00")
	}
}

func TestPlace_IsOpenAt_Closed(t *testing.T) {
	p := &domain.Place{
		Name: "Cafe",
		OpeningHours: &domain.OpeningHours{
			Periods: []domain.Period{
				{DayOfWeek: 1, OpenTime: "08:00", CloseTime: "20:00"},
			},
		},
	}
	if p.IsOpenAt(1, "21:00") {
		t.Error("expected place to be closed on Monday at 21:00")
	}
}

func TestPlace_IsOpenAt_WrongDay(t *testing.T) {
	p := &domain.Place{
		Name: "Cafe",
		OpeningHours: &domain.OpeningHours{
			Periods: []domain.Period{
				{DayOfWeek: 1, OpenTime: "08:00", CloseTime: "20:00"},
			},
		},
	}
	if p.IsOpenAt(0, "12:00") {
		t.Error("expected place to be closed on Sunday")
	}
}
