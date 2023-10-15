package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/tonytcb/flight-path-tracker/pkg/domain"
)

func TestFlightTracker_Track(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		flights domain.Flights
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Flight
		wantErr bool
	}{
		{
			name: "should successfully track a small flight list",
			args: args{
				ctx: context.Background(),
				flights: []*domain.Flight{
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
			},
			want: &domain.Flight{
				Source:      "SFO",
				Destination: "EWR",
			},
			wantErr: false,
		},
		{
			name: "should successfully track a long flight list",
			args: args{
				ctx:     context.Background(),
				flights: longListOfBrazilianFlights(),
			},
			want: &domain.Flight{
				Source:      "CWB",
				Destination: "VIX",
			},
			wantErr: false,
		},
		{
			name: "should error on invalid flight list",
			args: args{
				ctx:     context.Background(),
				flights: []*domain.Flight{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			f := NewFlightTracker()

			got, err := f.Track(tt.args.ctx, tt.args.flights)
			if (err != nil) != tt.wantErr {
				t.Errorf("Track() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Track() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// source: https://getbybus.com/en/blog/airports-brazil/
func longListOfBrazilianFlights() domain.Flights {
	return []*domain.Flight{
		{
			Source:      "MGA",
			Destination: "GRU",
		},
		{
			Source:      "REC",
			Destination: "CGH",
		},
		{
			Source:      "BSB",
			Destination: "MGA",
		},
		{
			Source:      "GSO",
			Destination: "SSA",
		},
		{
			Source:      "CNF",
			Destination: "ATL",
		},
		{
			Source:      "CGH",
			Destination: "VCP",
		},
		{
			Source:      "GRU",
			Destination: "POA",
		},
		{
			Source:      "SSA",
			Destination: "FOR",
		},
		{
			Source:      "CWB", // first origin
			Destination: "BEL",
		},
		{
			Source:      "FLN",
			Destination: "VIX", // final destination
		},
		{
			Source:      "NAT",
			Destination: "BSB",
		},
		{
			Source:      "VCP",
			Destination: "CGB",
		},
		{
			Source:      "BEL",
			Destination: "NAT",
		},
		{
			Source:      "CGB",
			Destination: "FLN",
		},
		{
			Source:      "ATL",
			Destination: "GSO",
		},
		{
			Source:      "MCZ",
			Destination: "REC",
		},
		{
			Source:      "POA",
			Destination: "IGU",
		},
		{
			Source:      "FOR",
			Destination: "MCZ",
		},
		{
			Source:      "IGU",
			Destination: "CNF",
		},
	}
}
