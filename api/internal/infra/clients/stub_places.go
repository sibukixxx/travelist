package clients

import (
	"context"
	"fmt"
	"strings"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

// StubPlacesClient is a development stub that returns hardcoded place data.
type StubPlacesClient struct {
	places map[string][]domain.Place
	byID   map[string]domain.Place
}

// NewStubPlacesClient returns a new StubPlacesClient with sample data.
func NewStubPlacesClient() *StubPlacesClient {
	kyoto := []domain.Place{
		{ID: "kyoto_kinkakuji", GooglePlaceID: "gp_kinkakuji", Name: "金閣寺", Lat: 35.0394, Lng: 135.7292, Types: []string{"temple", "tourist_attraction"}, PriceLevel: 1, Rating: 4.6, Address: "京都市北区金閣寺町1"},
		{ID: "kyoto_fushimi", GooglePlaceID: "gp_fushimi", Name: "伏見稲荷大社", Lat: 34.9671, Lng: 135.7727, Types: []string{"shrine", "tourist_attraction"}, PriceLevel: 0, Rating: 4.7, Address: "京都市伏見区深草薮之内町68"},
		{ID: "kyoto_arashiyama", GooglePlaceID: "gp_arashiyama", Name: "嵐山竹林", Lat: 35.0094, Lng: 135.6722, Types: []string{"park", "tourist_attraction"}, PriceLevel: 0, Rating: 4.5, Address: "京都市右京区嵯峨天龍寺芒ノ馬場町"},
		{ID: "kyoto_kiyomizu", GooglePlaceID: "gp_kiyomizu", Name: "清水寺", Lat: 34.9949, Lng: 135.7850, Types: []string{"temple", "tourist_attraction"}, PriceLevel: 1, Rating: 4.5, Address: "京都市東山区清水1丁目294"},
		{ID: "kyoto_nishiki", GooglePlaceID: "gp_nishiki", Name: "錦市場", Lat: 35.0050, Lng: 135.7649, Types: []string{"market", "food"}, PriceLevel: 2, Rating: 4.3, Address: "京都市中京区錦小路通"},
	}

	tokyo := []domain.Place{
		{ID: "tokyo_sensoji", GooglePlaceID: "gp_sensoji", Name: "浅草寺", Lat: 35.7148, Lng: 139.7967, Types: []string{"temple", "tourist_attraction"}, PriceLevel: 0, Rating: 4.5, Address: "東京都台東区浅草2丁目3-1"},
		{ID: "tokyo_meiji", GooglePlaceID: "gp_meiji", Name: "明治神宮", Lat: 35.6764, Lng: 139.6993, Types: []string{"shrine", "tourist_attraction"}, PriceLevel: 0, Rating: 4.6, Address: "東京都渋谷区代々木神園町1-1"},
		{ID: "tokyo_shibuya", GooglePlaceID: "gp_shibuya", Name: "渋谷スクランブル交差点", Lat: 35.6595, Lng: 139.7004, Types: []string{"landmark", "tourist_attraction"}, PriceLevel: 0, Rating: 4.3, Address: "東京都渋谷区道玄坂2丁目"},
		{ID: "tokyo_ueno", GooglePlaceID: "gp_ueno", Name: "上野公園", Lat: 35.7146, Lng: 139.7734, Types: []string{"park", "tourist_attraction"}, PriceLevel: 0, Rating: 4.4, Address: "東京都台東区上野公園"},
		{ID: "tokyo_tsukiji", GooglePlaceID: "gp_tsukiji", Name: "築地場外市場", Lat: 35.6654, Lng: 139.7707, Types: []string{"market", "food"}, PriceLevel: 2, Rating: 4.3, Address: "東京都中央区築地4丁目"},
	}

	c := &StubPlacesClient{
		places: map[string][]domain.Place{
			"京都": kyoto,
			"東京": tokyo,
		},
		byID: make(map[string]domain.Place),
	}

	for _, list := range c.places {
		for _, p := range list {
			c.byID[p.ID] = p
		}
	}

	return c
}

func (c *StubPlacesClient) SearchPlaces(_ context.Context, query string, _, _ float64) ([]domain.Place, error) {
	// Exact match
	if places, ok := c.places[query]; ok {
		return places, nil
	}
	// Partial match
	for key, places := range c.places {
		if strings.Contains(query, key) || strings.Contains(key, query) {
			return places, nil
		}
	}
	// Fallback: return Kyoto data for any unknown query
	return c.places["京都"], nil
}

func (c *StubPlacesClient) GetPlaceDetails(_ context.Context, placeID string) (*domain.Place, error) {
	p, ok := c.byID[placeID]
	if !ok {
		return nil, fmt.Errorf("place %q not found", placeID)
	}
	return &p, nil
}
