export interface Place {
  id: string
  google_place_id: string
  name: string
  lat: number
  lng: number
  types: string[]
  opening_hours?: OpeningHours
  price_level: number
  rating: number
  address: string
}

export interface OpeningHours {
  periods: Period[]
}

export interface Period {
  day_of_week: number
  open_time: string
  close_time: string
}

export interface Itinerary {
  id: string
  title: string
  destination: string
  start_date: string
  end_date: string
  days: DayPlan[]
  created_at: string
  updated_at: string
}

export interface DayPlan {
  day_number: number
  date: string
  activities: Activity[]
}

export interface Activity {
  order: number
  place_id: string
  place?: Place
  start_time: string
  end_time: string
  duration_min: number
  travel_from_prev?: TravelSegment
  note?: string
  estimated_cost_yen: number
}

export interface TravelSegment {
  mode: 'walk' | 'train' | 'bus' | 'taxi' | 'driving'
  duration_min: number
  distance_m: number
  estimated_cost_yen: number
}

export interface Violation {
  type: string
  day_number: number
  activity_idx: number
  message: string
}

export interface PlanRequest {
  destination: string
  num_days: number
  start_date: string
  preferences: {
    interests: string[]
    budget: 'budget' | 'moderate' | 'luxury'
    travel_style: 'relaxed' | 'active' | 'balanced'
    total_budget_yen?: number
  }
  constraint: {
    max_walk_distance_m: number
    max_activities_day: number
    earliest_start_time: string
    latest_end_time: string
  }
}

export interface DayCost {
  day_number: number
  cost_yen: number
}

export interface BudgetSummary {
  total_cost_yen: number
  daily_costs: DayCost[]
}

export interface GenerateResult {
  itinerary: Itinerary
  violations: Violation[]
  budget_summary: BudgetSummary
}
