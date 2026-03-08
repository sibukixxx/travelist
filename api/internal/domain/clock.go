package domain

import "time"

// Clock is an interface for getting the current time (testable).
type Clock interface {
	Now() time.Time
}
