package flightparser

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/tonytcb/flight-path-tracker/pkg/domain"
)

// JSONParser implements an improved version of the payload provided in the examples, considering a list of objects
type JSONParser struct {
}

func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

func (p *JSONParser) Parse(ctx context.Context, raw []byte) (domain.Flights, error) {
	type flight struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	var payload []*flight
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, errors.Wrap(err, "error to json decode payload")
	}

	var output = make([]*domain.Flight, 0)
	for k, v := range payload {
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), "context done while parsing payload")
		default:
		}

		if v.Source == "" {
			return nil, errors.Errorf("source value can not be empty on flight number %d", k)
		}

		if v.Destination == "" {
			return nil, errors.Errorf("destination value can not be empty on flight number %d", k)
		}

		output = append(output, domain.NewFlight(domain.Airport(v.Source), domain.Airport(v.Destination)))
	}

	return output, nil
}
