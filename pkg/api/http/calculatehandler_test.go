package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"go.uber.org/mock/gomock"

	"github.com/tonytcb/flight-path-tracker/pkg/domain"
)

func TestFlightCalculatorHandler_Handle(t *testing.T) {
	t.Parallel()

	var (
		rawBody1 = `[{"source":"IND","destination":"EWR"},{"source":"SFO","destination":"ATL"},{"source":"GSO","destination":"IND"},{"source":"ATL","destination":"GSO"}]`
		flights1 = []*domain.Flight{
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
		}
		flight1 = domain.NewFlight("SFO", "EWR")
	)

	type fields struct {
		parser  func(*gomock.Controller) FlightsParser
		tracker func(*gomock.Controller) FlightsTracker
	}
	type args struct {
		responseWriter http.ResponseWriter
		request        *http.Request
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name: "should calculate a flight path successfully",
			fields: fields{
				parser: func(ctrl *gomock.Controller) FlightsParser {
					parserMock := NewMockFlightsParser(ctrl)
					parserMock.EXPECT().
						Parse(gomock.Any(), []byte(rawBody1)).
						Return(flights1, nil).
						Times(1)

					return parserMock
				},
				tracker: func(ctrl *gomock.Controller) FlightsTracker {
					trackerMock := NewMockFlightsTracker(ctrl)
					trackerMock.EXPECT().
						Track(gomock.Any(), flights1).
						Return(flight1, nil).
						Times(1)

					return trackerMock
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        newRequest(t, "localhost:8080", http.MethodPost, rawBody1),
			},
			wantStatusCode:   200,
			wantResponseBody: `{"source":"SFO","destination":"EWR"}`,
		},

		{
			name: "should error on invalid http method",
			fields: fields{
				parser: func(ctrl *gomock.Controller) FlightsParser {
					return nil
				},
				tracker: func(ctrl *gomock.Controller) FlightsTracker {
					return nil
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        newRequest(t, "localhost:8080", http.MethodPut, rawBody1),
			},
			wantStatusCode:   405,
			wantResponseBody: ``,
		},
		{
			name: "should error on parser",
			fields: fields{
				parser: func(ctrl *gomock.Controller) FlightsParser {
					parserMock := NewMockFlightsParser(ctrl)
					parserMock.EXPECT().
						Parse(gomock.Any(), []byte(rawBody1)).
						Return(nil, errors.New("invalid json")).
						Times(1)

					return parserMock
				},
				tracker: func(ctrl *gomock.Controller) FlightsTracker {
					return nil
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        newRequest(t, "localhost:8080", http.MethodPost, rawBody1),
			},
			wantStatusCode:   400,
			wantResponseBody: `{"error":"error to parse json body: invalid json"}`,
		},
		{
			name: "should error on tracker",
			fields: fields{
				parser: func(ctrl *gomock.Controller) FlightsParser {
					parserMock := NewMockFlightsParser(ctrl)
					parserMock.EXPECT().
						Parse(gomock.Any(), []byte(rawBody1)).
						Return(flights1, nil).
						Times(1)

					return parserMock
				},
				tracker: func(ctrl *gomock.Controller) FlightsTracker {
					trackerMock := NewMockFlightsTracker(ctrl)
					trackerMock.EXPECT().
						Track(gomock.Any(), flights1).
						Return(nil, domain.ErrInvalidItinerary).
						Times(1)

					return trackerMock
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        newRequest(t, "localhost:8080", http.MethodPost, rawBody1),
			},
			wantStatusCode:   422,
			wantResponseBody: `{"error":"error to calculate original flight: invalid itinerary data"}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			h := &FlightCalculatorHandler{
				parser:  tt.fields.parser(mockCtrl),
				tracker: tt.fields.tracker(mockCtrl),
			}
			h.Handle(tt.args.responseWriter, tt.args.request)

			if v, ok := tt.args.responseWriter.(*httptest.ResponseRecorder); ok {
				httpResponse := v.Result()
				defer httpResponse.Body.Close()

				assertHTTPResponse(t, httpResponse, tt.wantStatusCode, tt.wantResponseBody)

				return
			}

			t.Error("impossible to assert http response due to type incompatibility")
		})
	}
}

func newRequest(t *testing.T, endpoint string, method string, payload string) *http.Request {
	body := strings.NewReader(payload)

	r, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		t.Fatalf("error to build request: %v", err)
	}

	return r
}
