package clock

import "time"

// Clock is an interface for getting the current time (testable).
type Clock interface {
	Now() time.Time
}

// RealClock uses the system clock.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// FixedClock returns a fixed time (for testing).
type FixedClock struct {
	Time time.Time
}

func (c FixedClock) Now() time.Time { return c.Time }
