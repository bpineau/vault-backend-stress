package metrics

import "time"

// Status represents a request status
type Status int

const (
	// Success represent a successful request
	Success Status = iota

	// Error represent a failed request
	Error
)

// Sink receive metrics
type Sink interface {
	Observe(Status, time.Duration)
	Dump() (*Point, error)
}

// Point is a snapshot of metrics state
type Point struct {
	SuccessRate  float64
	ErrorsRate   float64
	SuccessCount float64
	ErrorsCount  float64
	P50          float64
	P95          float64
	P99          float64
	Date         time.Time
}
