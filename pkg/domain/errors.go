package domain

import "errors"

var (
	ErrEmptyFlightsList = errors.New("there are no flights")
	ErrInvalidItinerary = errors.New("invalid itinerary data")
)
