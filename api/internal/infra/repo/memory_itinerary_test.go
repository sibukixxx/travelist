package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
)

func newTestItinerary(id string) *domain.Itinerary {
	return &domain.Itinerary{
		ID:          id,
		Title:       "Test Trip",
		Destination: "京都",
		StartDate:   time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC),
		Days: []domain.DayPlan{
			{
				DayNumber: 1,
				Date:      time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
				Activities: []domain.Activity{
					{Order: 0, PlaceID: "p1", StartTime: "09:00", EndTime: "11:00", DurationMin: 120},
				},
			},
		},
		CreatedAt: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestMemoryItineraryRepositorySave(t *testing.T) {
	t.Run("saves and retrieves itinerary by ID", func(t *testing.T) {
		r := repo.NewMemoryItineraryRepository()
		ctx := context.Background()
		it := newTestItinerary("itn_1")

		if err := r.Save(ctx, it); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := r.FindByID(ctx, "itn_1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if diff := cmp.Diff(it, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("overwrites existing itinerary with same ID", func(t *testing.T) {
		r := repo.NewMemoryItineraryRepository()
		ctx := context.Background()

		it1 := newTestItinerary("itn_1")
		it1.Title = "Original"
		_ = r.Save(ctx, it1)

		it2 := newTestItinerary("itn_1")
		it2.Title = "Updated"
		_ = r.Save(ctx, it2)

		got, _ := r.FindByID(ctx, "itn_1")
		if got.Title != "Updated" {
			t.Errorf("Title = %q, want %q", got.Title, "Updated")
		}
	})
}

func TestMemoryItineraryRepositoryFindByID(t *testing.T) {
	t.Run("returns error when ID does not exist", func(t *testing.T) {
		r := repo.NewMemoryItineraryRepository()
		_, err := r.FindByID(context.Background(), "non-existent")
		if err == nil {
			t.Error("expected error for non-existent ID, got nil")
		}
	})
}

func TestMemoryItineraryRepositoryList(t *testing.T) {
	t.Run("returns all saved itineraries", func(t *testing.T) {
		r := repo.NewMemoryItineraryRepository()
		ctx := context.Background()

		_ = r.Save(ctx, newTestItinerary("itn_1"))
		_ = r.Save(ctx, newTestItinerary("itn_2"))
		_ = r.Save(ctx, newTestItinerary("itn_3"))

		list, err := r.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list) != 3 {
			t.Fatalf("expected 3 items, got %d", len(list))
		}

		ids := make([]string, len(list))
		for i, it := range list {
			ids[i] = it.ID
		}
		wantIDs := []string{"itn_1", "itn_2", "itn_3"}
		if diff := cmp.Diff(wantIDs, ids, cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
			t.Errorf("IDs mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("returns empty slice when no itineraries exist", func(t *testing.T) {
		r := repo.NewMemoryItineraryRepository()
		list, err := r.List(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list) != 0 {
			t.Errorf("expected empty list, got %d items", len(list))
		}
	})
}

func TestMemoryItineraryRepositoryDelete(t *testing.T) {
	t.Run("removes itinerary by ID", func(t *testing.T) {
		r := repo.NewMemoryItineraryRepository()
		ctx := context.Background()

		_ = r.Save(ctx, newTestItinerary("itn_1"))
		if err := r.Delete(ctx, "itn_1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err := r.FindByID(ctx, "itn_1")
		if err == nil {
			t.Error("expected error after delete, got nil")
		}
	})

	t.Run("returns error when deleting non-existent ID", func(t *testing.T) {
		r := repo.NewMemoryItineraryRepository()
		err := r.Delete(context.Background(), "non-existent")
		if err == nil {
			t.Error("expected error for non-existent ID, got nil")
		}
	})
}

// Verify interface compliance at compile time.
var _ repo.ItineraryRepository = (*repo.MemoryItineraryRepository)(nil)
