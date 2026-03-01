package clients_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sibukixxx/travelist/api/internal/infra/clients"
)

func TestStubPlacesClientSearchPlaces(t *testing.T) {
	t.Run("returns places matching destination keyword", func(t *testing.T) {
		client := clients.NewStubPlacesClient()
		places, err := client.SearchPlaces(context.Background(), "京都", 0, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(places) == 0 {
			t.Fatal("expected places, got empty slice")
		}
		for _, p := range places {
			if p.ID == "" {
				t.Error("place ID should not be empty")
			}
			if p.Name == "" {
				t.Error("place Name should not be empty")
			}
		}
	})

	t.Run("returns fallback places for unknown destination", func(t *testing.T) {
		client := clients.NewStubPlacesClient()
		places, err := client.SearchPlaces(context.Background(), "unknown-city", 0, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(places) == 0 {
			t.Fatal("expected fallback places, got empty slice")
		}
	})

	t.Run("each place has non-zero lat/lng", func(t *testing.T) {
		client := clients.NewStubPlacesClient()
		places, err := client.SearchPlaces(context.Background(), "京都", 0, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, p := range places {
			if p.Lat == 0 || p.Lng == 0 {
				t.Errorf("place %q has zero lat/lng", p.Name)
			}
		}
	})
}

func TestStubPlacesClientGetPlaceDetails(t *testing.T) {
	t.Run("returns place when ID exists in stub data", func(t *testing.T) {
		client := clients.NewStubPlacesClient()
		// First search to populate known places
		places, _ := client.SearchPlaces(context.Background(), "京都", 0, 0)
		if len(places) == 0 {
			t.Fatal("need at least one place for this test")
		}

		target := places[0]
		got, err := client.GetPlaceDetails(context.Background(), target.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if diff := cmp.Diff(target.Name, got.Name); diff != "" {
			t.Errorf("Name mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("returns error when ID does not exist", func(t *testing.T) {
		client := clients.NewStubPlacesClient()
		_, err := client.GetPlaceDetails(context.Background(), "non-existent-id")
		if err == nil {
			t.Error("expected error for non-existent ID, got nil")
		}
	})
}

// Verify interface compliance at compile time.
var _ clients.PlacesClient = (*clients.StubPlacesClient)(nil)
