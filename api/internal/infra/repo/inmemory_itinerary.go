package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

// InMemoryItineraryRepository is a thread-safe in-memory implementation of ItineraryRepository.
type InMemoryItineraryRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Itinerary
}

// NewInMemoryItineraryRepository returns a new InMemoryItineraryRepository.
func NewInMemoryItineraryRepository() *InMemoryItineraryRepository {
	return &InMemoryItineraryRepository{
		data: make(map[string]*domain.Itinerary),
	}
}

func (r *InMemoryItineraryRepository) Save(_ context.Context, itinerary *domain.Itinerary) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[itinerary.ID] = itinerary
	return nil
}

func (r *InMemoryItineraryRepository) FindByID(_ context.Context, id string) (*domain.Itinerary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	it, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("itinerary %q not found", id)
	}
	return it, nil
}

func (r *InMemoryItineraryRepository) List(_ context.Context) ([]*domain.Itinerary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Itinerary, 0, len(r.data))
	for _, it := range r.data {
		result = append(result, it)
	}
	return result, nil
}

func (r *InMemoryItineraryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("itinerary %q not found", id)
	}
	delete(r.data, id)
	return nil
}
