package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/tonytcb/flight-path-tracker/pkg/domain"
)

const (
	timeoutDefault = 10 * time.Second
)

type FlightsTracker interface {
	Track(context.Context, domain.Flights) (*domain.Flight, error)
}

type FlightsParser interface {
	Parse(context.Context, []byte) (domain.Flights, error)
}

type FlightCalculatorHandler struct {
	parser  FlightsParser
	tracker FlightsTracker
}

func NewFlightCalculatorHandler(
	parser FlightsParser,
	tracker FlightsTracker,
) *FlightCalculatorHandler {
	return &FlightCalculatorHandler{parser: parser, tracker: tracker}
}

func (h *FlightCalculatorHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeoutDefault)
	defer cancel()

	var output = jsonOutput{w: w}

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		_ = output.internalServerError(err, "error to read body")
		return
	}
	defer r.Body.Close()

	flights, err := h.parser.Parse(ctx, rawBody)
	if err != nil {
		_ = output.badRequest(err, "error to parse json body")
		return
	}

	flightResponse, err := h.tracker.Track(ctx, flights)
	if err != nil {
		_ = output.domainError(err, "error to calculate original flight")
		return
	}

	_ = output.ok(flightResponse)
}
