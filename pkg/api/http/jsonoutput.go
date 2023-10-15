package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tonytcb/flight-path-tracker/pkg/domain"
)

type httpError struct {
	Error string `json:"error"`
}

type jsonOutput struct {
	w http.ResponseWriter
}

func (o jsonOutput) ok(flight *domain.Flight) error {
	output := struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}{
		Source:      string(flight.Source),
		Destination: string(flight.Destination),
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return errors.Wrap(err, "error to encode flight output")
	}

	o.w.Header().Add("Content-Type", "application/json")
	o.w.WriteHeader(http.StatusOK)
	_, err = o.w.Write(bytes)

	return errors.Wrap(err, "error to write response")
}

func (o jsonOutput) internalServerError(err error, details string) error {
	output := httpError{
		Error: fmt.Sprintf("%s: %s", details, err.Error()),
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return errors.Wrap(err, "error to encode error output")
	}

	o.w.Header().Add("Content-Type", "application/json")
	o.w.WriteHeader(http.StatusInternalServerError)
	_, err = o.w.Write(bytes)

	return errors.Wrap(err, "error to write response")
}

func (o jsonOutput) badRequest(err error, details string) error {
	output := httpError{
		Error: fmt.Sprintf("%s: %s", details, err.Error()),
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return errors.Wrap(err, "error to encode error output")
	}

	o.w.Header().Add("Content-Type", "application/json")
	o.w.WriteHeader(http.StatusBadRequest)
	_, err = o.w.Write(bytes)

	return errors.Wrap(err, "error to write response")
}

func (o jsonOutput) domainError(rootErr error, details string) error {
	output := httpError{
		Error: fmt.Sprintf("%s: %s", details, rootErr.Error()),
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return errors.Wrap(err, "error to encode error output")
	}

	o.w.Header().Add("Content-Type", "application/json")
	o.w.WriteHeader(translateDomainErr(rootErr))
	_, err = o.w.Write(bytes)

	return errors.Wrap(err, "error to write response")
}

func translateDomainErr(err error) int {
	switch {
	case errors.Is(err, domain.ErrEmptyFlightsList):
		return http.StatusUnprocessableEntity

	case errors.Is(err, domain.ErrInvalidItinerary):
		return http.StatusUnprocessableEntity

	default:
		return http.StatusServiceUnavailable
	}
}
