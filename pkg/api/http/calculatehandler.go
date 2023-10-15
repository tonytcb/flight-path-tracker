package http

import (
	"context"
	"fmt"
	"net/http"
)

type FlightCalculator interface {
	Parse(context.Context) error
}

type FlightCalculatorHandler struct {
	calculator FlightCalculator
}

func NewFlightCalculator(
	calculator FlightCalculator,
) *FlightCalculatorHandler {
	return &FlightCalculatorHandler{calculator: calculator}
}

func (h *FlightCalculatorHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	//w.Write([]byte("wip"))
	fmt.Fprintln(w, "wip")
}
