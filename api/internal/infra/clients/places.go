package clients

import (
	"context"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

// PlacesClient defines the interface for searching places.
type PlacesClient interface {
	SearchPlaces(ctx context.Context, query string, lat, lng float64) ([]domain.Place, error)
	GetPlaceDetails(ctx context.Context, placeID string) (*domain.Place, error)
}
