package repo

import (
	"context"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

// ItineraryRepository defines the persistence interface for itineraries.
type ItineraryRepository interface {
	Save(ctx context.Context, itinerary *domain.Itinerary) error
	FindByID(ctx context.Context, id string) (*domain.Itinerary, error)
	List(ctx context.Context) ([]*domain.Itinerary, error)
	Delete(ctx context.Context, id string) error
}
