package clients

import (
	"context"
	"fmt"
	"strings"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

// StubPlacesClient is an in-memory stub implementation of PlacesClient.
type StubPlacesClient struct {
	allPlaces map[string][]domain.Place
}

// NewStubPlacesClient creates a StubPlacesClient with hardcoded data.
func NewStubPlacesClient() *StubPlacesClient {
	data := map[string][]domain.Place{
		"京都": {
			{ID: "kyoto-1", Name: "金閣寺", Lat: 35.0394, Lng: 135.7292, Types: []string{"temple", "culture"}, Rating: 4.6, Address: "京都市北区金閣寺町1"},
			{ID: "kyoto-2", Name: "伏見稲荷大社", Lat: 34.9671, Lng: 135.7727, Types: []string{"shrine", "culture"}, Rating: 4.5, Address: "京都市伏見区深草藪之内町68"},
			{ID: "kyoto-3", Name: "清水寺", Lat: 34.9948, Lng: 135.7850, Types: []string{"temple", "culture"}, Rating: 4.5, Address: "京都市東山区清水1丁目294"},
			{ID: "kyoto-4", Name: "嵐山竹林", Lat: 35.0170, Lng: 135.6713, Types: []string{"nature", "scenic"}, Rating: 4.4, Address: "京都市右京区嵯峨天龍寺芒ノ馬場町"},
			{ID: "kyoto-5", Name: "錦市場", Lat: 35.0050, Lng: 135.7649, Types: []string{"food", "market"}, Rating: 4.3, Address: "京都市中京区錦小路通"},
		},
		"東京": {
			{ID: "tokyo-1", Name: "東京タワー", Lat: 35.6586, Lng: 139.7454, Types: []string{"landmark"}, Rating: 4.3, Address: "東京都港区芝公園4丁目2-8"},
			{ID: "tokyo-2", Name: "浅草寺", Lat: 35.7148, Lng: 139.7967, Types: []string{"temple", "culture"}, Rating: 4.5, Address: "東京都台東区浅草2丁目3-1"},
			{ID: "tokyo-3", Name: "明治神宮", Lat: 35.6764, Lng: 139.6993, Types: []string{"shrine", "nature"}, Rating: 4.5, Address: "東京都渋谷区代々木神園町1-1"},
		},
		"大阪": {
			{ID: "osaka-1", Name: "大阪城", Lat: 34.6873, Lng: 135.5262, Types: []string{"castle", "culture"}, Rating: 4.3, Address: "大阪市中央区大阪城1-1"},
			{ID: "osaka-2", Name: "道頓堀", Lat: 34.6687, Lng: 135.5013, Types: []string{"food", "entertainment"}, Rating: 4.4, Address: "大阪市中央区道頓堀"},
		},
	}

	// Default fallback places
	data["_default"] = []domain.Place{
		{ID: "default-1", Name: "観光スポットA", Lat: 35.6812, Lng: 139.7671, Types: []string{"sightseeing"}, Rating: 4.0, Address: "観光地A"},
		{ID: "default-2", Name: "観光スポットB", Lat: 35.6822, Lng: 139.7681, Types: []string{"sightseeing"}, Rating: 4.0, Address: "観光地B"},
	}

	return &StubPlacesClient{allPlaces: data}
}

func (s *StubPlacesClient) SearchPlaces(_ context.Context, query string, _, _ float64) ([]domain.Place, error) {
	// Try exact match first
	if places, ok := s.allPlaces[query]; ok {
		return places, nil
	}
	// Try prefix/contains match
	for key, places := range s.allPlaces {
		if key == "_default" {
			continue
		}
		if strings.Contains(key, query) || strings.Contains(query, key) {
			return places, nil
		}
	}
	return s.allPlaces["_default"], nil
}

func (s *StubPlacesClient) GetPlaceDetails(_ context.Context, placeID string) (*domain.Place, error) {
	for _, places := range s.allPlaces {
		for _, p := range places {
			if p.ID == placeID {
				return &p, nil
			}
		}
	}
	return nil, fmt.Errorf("place not found: %s", placeID)
}
