package usecase

import (
	"context"

	"github.com/pkg/errors"

	"github.com/tonytcb/flight-path-tracker/pkg/domain"
)

type FlightTracker struct {
}

func NewFlightTracker() *FlightTracker {
	return &FlightTracker{}
}

func (f *FlightTracker) Track(_ context.Context, flights domain.Flights) (*domain.Flight, error) {
	originalFlight, err := flights.OriginalSourceAndDestination()
	if err != nil {
		return originalFlight, errors.Wrap(err, "error to track flight")
	}

	return originalFlight, nil
}
