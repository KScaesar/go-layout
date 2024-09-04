package utility

import (
	"time"
)

// MockTimeNow
//
// Parameters:
// RFC3339 (string): The time value in RFC3339 format that will be used for simulation.
//
// Example usage:
//
//	MockTimeNow("2023-08-19T12:00:00Z")
//	MockTimeNow("2023-08-19T20:00:00+08:00")
func MockTimeNow(RFC3339 string) func() time.Time {
	t, err := time.Parse(time.RFC3339, RFC3339)
	if err != nil {
		panic(err)
	}
	return func() time.Time { return t }
}
