package domain

import "context"

// ItineraryRepository defines the persistence interface for itineraries.
type ItineraryRepository interface {
	Save(ctx context.Context, itinerary *Itinerary) error
	FindByID(ctx context.Context, id string) (*Itinerary, error)
	List(ctx context.Context) ([]*Itinerary, error)
	Delete(ctx context.Context, id string) error
}
