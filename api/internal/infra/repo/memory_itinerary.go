package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

// MemoryItineraryRepository is an in-memory implementation of ItineraryRepository.
type MemoryItineraryRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Itinerary
}

// NewMemoryItineraryRepository creates a new MemoryItineraryRepository.
func NewMemoryItineraryRepository() *MemoryItineraryRepository {
	return &MemoryItineraryRepository{
		data: make(map[string]*domain.Itinerary),
	}
}

func (r *MemoryItineraryRepository) Save(_ context.Context, itinerary *domain.Itinerary) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[itinerary.ID] = itinerary
	return nil
}

func (r *MemoryItineraryRepository) FindByID(_ context.Context, id string) (*domain.Itinerary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	it, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("itinerary not found: %s", id)
	}
	return it, nil
}

func (r *MemoryItineraryRepository) List(_ context.Context) ([]*domain.Itinerary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Itinerary, 0, len(r.data))
	for _, it := range r.data {
		result = append(result, it)
	}
	return result, nil
}

func (r *MemoryItineraryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("itinerary not found: %s", id)
	}
	delete(r.data, id)
	return nil
}
