package repo_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
)

func newTestItinerary(id string) *domain.Itinerary {
	now := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	return &domain.Itinerary{
		ID:          id,
		Title:       "京都 2日間の旅",
		Destination: "京都",
		StartDate:   now,
		EndDate:     now.AddDate(0, 0, 1),
		Days:        []domain.DayPlan{{DayNumber: 1, Date: now}},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestInMemoryItineraryRepository(t *testing.T) {
	t.Run("Save then FindByID returns the saved itinerary", func(t *testing.T) {
		r := repo.NewInMemoryItineraryRepository()
		ctx := context.Background()
		it := newTestItinerary("itn_1")

		if err := r.Save(ctx, it); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		got, err := r.FindByID(ctx, "itn_1")
		if err != nil {
			t.Fatalf("FindByID failed: %v", err)
		}
		if diff := cmp.Diff(it, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("FindByID returns error for non-existent ID", func(t *testing.T) {
		r := repo.NewInMemoryItineraryRepository()
		ctx := context.Background()

		_, err := r.FindByID(ctx, "non_existent")
		if err == nil {
			t.Fatal("expected error for non-existent ID, got nil")
		}
	})

	t.Run("List returns all saved itineraries", func(t *testing.T) {
		r := repo.NewInMemoryItineraryRepository()
		ctx := context.Background()

		it1 := newTestItinerary("itn_1")
		it2 := newTestItinerary("itn_2")
		_ = r.Save(ctx, it1)
		_ = r.Save(ctx, it2)

		got, err := r.List(ctx)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("expected 2 itineraries, got %d", len(got))
		}
	})

	t.Run("Delete removes itinerary so FindByID returns error", func(t *testing.T) {
		r := repo.NewInMemoryItineraryRepository()
		ctx := context.Background()

		it := newTestItinerary("itn_1")
		_ = r.Save(ctx, it)

		if err := r.Delete(ctx, "itn_1"); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err := r.FindByID(ctx, "itn_1")
		if err == nil {
			t.Fatal("expected error after Delete, got nil")
		}
	})

	t.Run("Delete returns error for non-existent ID", func(t *testing.T) {
		r := repo.NewInMemoryItineraryRepository()
		ctx := context.Background()

		err := r.Delete(ctx, "non_existent")
		if err == nil {
			t.Fatal("expected error for deleting non-existent ID, got nil")
		}
	})

	t.Run("Save overwrites existing itinerary with same ID", func(t *testing.T) {
		r := repo.NewInMemoryItineraryRepository()
		ctx := context.Background()

		it := newTestItinerary("itn_1")
		_ = r.Save(ctx, it)

		updated := newTestItinerary("itn_1")
		updated.Title = "更新された旅程"
		_ = r.Save(ctx, updated)

		got, err := r.FindByID(ctx, "itn_1")
		if err != nil {
			t.Fatalf("FindByID failed: %v", err)
		}
		if got.Title != "更新された旅程" {
			t.Errorf("expected updated title, got %q", got.Title)
		}
	})

	t.Run("concurrent access is safe", func(t *testing.T) {
		r := repo.NewInMemoryItineraryRepository()
		ctx := context.Background()

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				it := newTestItinerary("itn_concurrent")
				_ = r.Save(ctx, it)
				_, _ = r.FindByID(ctx, "itn_concurrent")
				_, _ = r.List(ctx)
			}(i)
		}
		wg.Wait()

		// If we get here without a race condition panic, the test passes.
	})
}
