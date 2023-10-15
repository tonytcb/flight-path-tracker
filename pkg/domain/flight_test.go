package domain

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestFlights_OriginalSourceAndDestination(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		flights Flights
		want    *Flight
		wantErr error
	}{
		{
			name:    "should error when there are no flights",
			flights: []*Flight{},
			want:    nil,
			wantErr: ErrEmptyFlightsList,
		},
		{
			name: "should error when the itinerary is invalid",
			flights: []*Flight{
				{
					Source:      "SFO",
					Destination: "EWR",
				},
				{
					Source:      "SFO",
					Destination: "EWR",
				},
			},
			want:    nil,
			wantErr: ErrInvalidItinerary,
		},
		{
			name: "should return EWR when there's only one flight",
			flights: []*Flight{
				{
					Source:      "SFO",
					Destination: "EWR",
				},
			},
			want: &Flight{
				Source:      "SFO",
				Destination: "EWR",
			},
		},
		{
			name: "should return EWR when there are 4 flights flight",
			flights: []*Flight{
				{
					Source:      "IND",
					Destination: "EWR",
				},
				{
					Source:      "SFO",
					Destination: "ATL",
				},
				{
					Source:      "GSO",
					Destination: "IND",
				},
				{
					Source:      "ATL",
					Destination: "GSO",
				},
			},
			want: &Flight{
				Source:      "SFO",
				Destination: "EWR",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.flights.OriginalSourceAndDestination()
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("OriginalSourceAndDestination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OriginalSourceAndDestination() got = %v, want %v", got, tt.want)
			}
		})
	}
}
