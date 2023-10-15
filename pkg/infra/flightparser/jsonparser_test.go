package flightparser

import (
	"context"
	"reflect"
	"testing"

	"github.com/tonytcb/flight-path-tracker/pkg/domain"
)

func TestJSONParser_Parse(t *testing.T) {
	t.Parallel()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type args struct {
		ctx context.Context
		raw []byte
	}
	tests := []struct {
		name    string
		args    args
		want    domain.Flights
		wantErr bool
	}{
		{
			name: "should parse a json containing two flights successfully ",
			args: args{
				ctx: context.Background(),
				raw: []byte(`[{"source":"IND","destination":"EWR"},{"source":"SFO","destination":"ATL"}]`),
			},
			want: []*domain.Flight{
				{
					Source:      "IND",
					Destination: "EWR",
				},
				{
					Source:      "SFO",
					Destination: "ATL",
				},
			},
			wantErr: false,
		},
		{
			name: "should parse an empty json successfully",
			args: args{
				ctx: context.Background(),
				raw: []byte(`[]`),
			},
			want:    []*domain.Flight{},
			wantErr: false,
		},
		{
			name: "should error on an invalid json",
			args: args{
				ctx: context.Background(),
				raw: []byte(`invalid json`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should error on context done",
			args: args{
				ctx: canceledCtx,
				raw: []byte(`[{"source":"IND","destination":"EWR"},{"source":"SFO","destination":"ATL"}]`),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			p := NewJSONParser()

			got, err := p.Parse(tt.args.ctx, tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
