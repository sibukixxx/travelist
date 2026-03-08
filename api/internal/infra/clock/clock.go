package clock

import "time"

// RealClock uses the system clock.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// FixedClock returns a fixed time (for testing).
type FixedClock struct {
	Time time.Time
}

func (c FixedClock) Now() time.Time { return c.Time }
