package clients_test

import (
	"context"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/infra/clients"
)

func TestStubPlacesClient(t *testing.T) {
	t.Run("SearchPlaces returns places matching query", func(t *testing.T) {
		c := clients.NewStubPlacesClient()
		ctx := context.Background()

		places, err := c.SearchPlaces(ctx, "京都", 0, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(places) == 0 {
			t.Fatal("expected at least one place, got 0")
		}
		for _, p := range places {
			if p.ID == "" || p.Name == "" {
				t.Errorf("place has empty ID or Name: %+v", p)
			}
		}
	})

	t.Run("SearchPlaces returns places for unknown query", func(t *testing.T) {
		c := clients.NewStubPlacesClient()
		ctx := context.Background()

		places, err := c.SearchPlaces(ctx, "unknown_destination", 0, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(places) == 0 {
			t.Fatal("expected fallback places for unknown query, got 0")
		}
	})

	t.Run("GetPlaceDetails returns place for known ID", func(t *testing.T) {
		c := clients.NewStubPlacesClient()
		ctx := context.Background()

		// First get a place via search to know a valid ID
		places, _ := c.SearchPlaces(ctx, "京都", 0, 0)
		if len(places) == 0 {
			t.Fatal("no places from search to test details")
		}

		place, err := c.GetPlaceDetails(ctx, places[0].ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if place.ID != places[0].ID {
			t.Errorf("expected ID %q, got %q", places[0].ID, place.ID)
		}
	})

	t.Run("GetPlaceDetails returns error for unknown ID", func(t *testing.T) {
		c := clients.NewStubPlacesClient()
		ctx := context.Background()

		_, err := c.GetPlaceDetails(ctx, "unknown_id")
		if err == nil {
			t.Fatal("expected error for unknown place ID, got nil")
		}
	})
}
