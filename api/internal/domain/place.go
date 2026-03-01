package domain

// Place represents a normalized location from Google Places API.
type Place struct {
	ID             string   `json:"id"`
	GooglePlaceID  string   `json:"google_place_id"`
	Name           string   `json:"name"`
	Lat            float64  `json:"lat"`
	Lng            float64  `json:"lng"`
	Types          []string `json:"types"`
	OpeningHours   *OpeningHours `json:"opening_hours,omitempty"`
	PriceLevel     int      `json:"price_level"`
	Rating         float64  `json:"rating"`
	Address        string   `json:"address"`
}

// OpeningHours represents normalized opening hours.
type OpeningHours struct {
	Periods []Period `json:"periods"`
}

// Period represents a single opening period.
type Period struct {
	DayOfWeek int    `json:"day_of_week"` // 0=Sunday, 6=Saturday
	OpenTime  string `json:"open_time"`   // "HH:MM"
	CloseTime string `json:"close_time"`  // "HH:MM"
}

// IsOpenAt checks if the place is open at the given day and time.
func (p *Place) IsOpenAt(dayOfWeek int, timeHHMM string) bool {
	if p.OpeningHours == nil {
		return true // assume open if no data
	}
	for _, period := range p.OpeningHours.Periods {
		if period.DayOfWeek == dayOfWeek &&
			period.OpenTime <= timeHHMM &&
			timeHHMM < period.CloseTime {
			return true
		}
	}
	return false
}
