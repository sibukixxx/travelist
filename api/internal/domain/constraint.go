package domain

import "fmt"

// Constraint represents a user-defined constraint for the itinerary.
type Constraint struct {
	MaxWalkDistanceM  int    `json:"max_walk_distance_m"`
	MaxActivitiesDay  int    `json:"max_activities_day"`
	EarliestStartTime string `json:"earliest_start_time"` // "HH:MM"
	LatestEndTime     string `json:"latest_end_time"`     // "HH:MM"
}

// DefaultConstraint returns sensible defaults.
func DefaultConstraint() Constraint {
	return Constraint{
		MaxWalkDistanceM:  2000,
		MaxActivitiesDay:  6,
		EarliestStartTime: "08:00",
		LatestEndTime:     "21:00",
	}
}

// ViolationType identifies the kind of constraint violation.
type ViolationType string

const (
	ViolationClosedPlace       ViolationType = "closed_place"
	ViolationExcessiveWalk     ViolationType = "excessive_walk"
	ViolationTimeOverlap       ViolationType = "time_overlap"
	ViolationTooManyActivities ViolationType = "too_many_activities"
	ViolationOutsideHours      ViolationType = "outside_hours"
	ViolationImpossibleTravel  ViolationType = "impossible_travel"
	ViolationBudgetExceeded    ViolationType = "budget_exceeded"
)

// Violation represents a detected problem in the itinerary.
type Violation struct {
	Type        ViolationType `json:"type"`
	DayNumber   int           `json:"day_number"`
	ActivityIdx int           `json:"activity_idx"`
	Message     string        `json:"message"`
}

// ValidateDayPlan checks a day plan against the given constraint and returns violations.
func ValidateDayPlan(day DayPlan, constraint Constraint) []Violation {
	var violations []Violation

	// Check number of activities
	if len(day.Activities) > constraint.MaxActivitiesDay {
		violations = append(violations, Violation{
			Type:      ViolationTooManyActivities,
			DayNumber: day.DayNumber,
			Message:   fmt.Sprintf("day %d has %d activities, max is %d", day.DayNumber, len(day.Activities), constraint.MaxActivitiesDay),
		})
	}

	for i, act := range day.Activities {
		// Check outside hours
		if act.StartTime < constraint.EarliestStartTime {
			violations = append(violations, Violation{
				Type:        ViolationOutsideHours,
				DayNumber:   day.DayNumber,
				ActivityIdx: i,
				Message:     fmt.Sprintf("activity %d starts at %s, before earliest %s", i, act.StartTime, constraint.EarliestStartTime),
			})
		}
		if act.EndTime > constraint.LatestEndTime {
			violations = append(violations, Violation{
				Type:        ViolationOutsideHours,
				DayNumber:   day.DayNumber,
				ActivityIdx: i,
				Message:     fmt.Sprintf("activity %d ends at %s, after latest %s", i, act.EndTime, constraint.LatestEndTime),
			})
		}

		// Check opening hours
		if act.Place != nil {
			dow := int(day.Date.Weekday())
			if !act.Place.IsOpenAt(dow, act.StartTime) {
				violations = append(violations, Violation{
					Type:        ViolationClosedPlace,
					DayNumber:   day.DayNumber,
					ActivityIdx: i,
					Message:     fmt.Sprintf("activity %d: %s is closed at %s", i, act.Place.Name, act.StartTime),
				})
			}
		}

		// Check excessive walk distance
		if act.TravelFromPrev != nil &&
			act.TravelFromPrev.Mode == TravelModeWalk &&
			act.TravelFromPrev.DistanceM > constraint.MaxWalkDistanceM {
			violations = append(violations, Violation{
				Type:        ViolationExcessiveWalk,
				DayNumber:   day.DayNumber,
				ActivityIdx: i,
				Message:     fmt.Sprintf("activity %d: walk distance %dm exceeds max %dm", i, act.TravelFromPrev.DistanceM, constraint.MaxWalkDistanceM),
			})
		}

		// Check time overlap with next activity
		if i > 0 {
			prev := day.Activities[i-1]
			if prev.EndTime > act.StartTime {
				violations = append(violations, Violation{
					Type:        ViolationTimeOverlap,
					DayNumber:   day.DayNumber,
					ActivityIdx: i,
					Message:     fmt.Sprintf("activity %d starts at %s but previous ends at %s", i, act.StartTime, prev.EndTime),
				})
			}
		}
	}

	return violations
}
