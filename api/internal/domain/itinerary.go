package domain

import "time"

// Itinerary represents a complete travel plan.
type Itinerary struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Destination string       `json:"destination"`
	StartDate   time.Time    `json:"start_date"`
	EndDate     time.Time    `json:"end_date"`
	Days        []DayPlan    `json:"days"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// DayPlan represents a single day's plan within an itinerary.
type DayPlan struct {
	DayNumber  int        `json:"day_number"`
	Date       time.Time  `json:"date"`
	Activities []Activity `json:"activities"`
}

// Activity represents a single activity in a day plan.
type Activity struct {
	Order         int    `json:"order"`
	PlaceID       string `json:"place_id"`
	Place         *Place `json:"place,omitempty"`
	StartTime     string `json:"start_time"`      // "HH:MM"
	EndTime       string `json:"end_time"`         // "HH:MM"
	DurationMin   int    `json:"duration_min"`
	TravelFromPrev *TravelSegment `json:"travel_from_prev,omitempty"`
	Note          string `json:"note,omitempty"`
}

// TravelSegment represents travel between two activities.
type TravelSegment struct {
	Mode        TravelMode `json:"mode"`
	DurationMin int        `json:"duration_min"`
	DistanceM   int        `json:"distance_m"`
}

// TravelMode represents a mode of transportation.
type TravelMode string

const (
	TravelModeWalk    TravelMode = "walk"
	TravelModeTrain   TravelMode = "train"
	TravelModeBus     TravelMode = "bus"
	TravelModeTaxi    TravelMode = "taxi"
	TravelModeDriving TravelMode = "driving"
)

// NumDays returns the number of days in the itinerary.
func (it *Itinerary) NumDays() int {
	return int(it.EndDate.Sub(it.StartDate).Hours()/24) + 1
}
