package domain

import "errors"

var (
	ErrEmptyFlightsList = errors.New("there are not flights")
	ErrInvalidItinerary = errors.New("invalid itinerary data")
)
