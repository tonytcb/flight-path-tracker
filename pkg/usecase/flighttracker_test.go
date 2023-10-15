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
			name: "should successfully track a flight list",
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
