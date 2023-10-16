package flightparser

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/tonytcb/flight-path-tracker/pkg/domain"
)

// JSONOfArraysParser implements exactly the same json provided in the examples
type JSONOfArraysParser struct {
}

func NewJSONOfArraysParser() *JSONOfArraysParser {
	return &JSONOfArraysParser{}
}

func (p *JSONOfArraysParser) Parse(_ context.Context, raw []byte) (domain.Flights, error) {
	var payload [][]string
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, errors.Wrap(err, "error to json decode payload")
	}

	const expectedPositions = 2

	var output = make([]*domain.Flight, 0)
	for k, v := range payload {
		if len(v) != expectedPositions {
			return nil, errors.Errorf("invalid flight %d, expected exactly 2 positions", k)
		}

		var (
			source      = v[0]
			destination = v[1]
		)

		if source == "" {
			return nil, errors.Errorf("source value can not be empty on flight number %d", k)
		}

		if destination == "" {
			return nil, errors.Errorf("destination value can not be empty on flight number %d", k)
		}

		output = append(output, domain.NewFlight(domain.Airport(source), domain.Airport(destination)))
	}

	return output, nil
}
