package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"

	"github.com/tonytcb/flight-path-tracker/pkg/domain"
)

func Test_jsonOutput_ok(t *testing.T) {
	t.Parallel()

	const (
		expectedStatusCode = 200
		expectedPayload    = `{"source":"SFO","destination":"EWR"}`
	)

	var (
		flight = &domain.Flight{
			Source:      "SFO",
			Destination: "EWR",
		}
		responseWriter = httptest.NewRecorder()
	)

	err := jsonOutput{w: responseWriter}.ok(flight)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assertHTTPResponse(t, responseWriter.Result(), expectedStatusCode, expectedPayload)
}

func Test_jsonOutput_internalServerError(t *testing.T) {
	t.Parallel()

	const (
		expectedStatusCode = 500
		expectedPayload    = `{"error":"error x: network error"}`
	)

	var responseWriter = httptest.NewRecorder()

	err := jsonOutput{w: responseWriter}.internalServerError(errors.New("network error"), "error x")
	if err != nil {
		t.Fatalf(err.Error())
	}

	assertHTTPResponse(t, responseWriter.Result(), expectedStatusCode, expectedPayload)
}

func Test_jsonOutput_badRequest(t *testing.T) {
	t.Parallel()

	const (
		expectedStatusCode = 400
		expectedPayload    = `{"error":"error to parse input: error on position x"}`
	)

	var responseWriter = httptest.NewRecorder()

	err := jsonOutput{w: responseWriter}.badRequest(errors.New("error on position x"), "error to parse input")
	if err != nil {
		t.Fatalf(err.Error())
	}

	assertHTTPResponse(t, responseWriter.Result(), expectedStatusCode, expectedPayload)
}

func Test_jsonOutput_domainError(t *testing.T) {
	t.Parallel()

	type fields struct {
		responseWriter *httptest.ResponseRecorder
	}
	type args struct {
		err     error
		details string
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name: "empty flight list",
			fields: fields{
				responseWriter: httptest.NewRecorder(),
			},
			args: args{
				err:     domain.ErrEmptyFlightsList,
				details: "error to calculate path",
			},
			wantStatusCode:   422,
			wantResponseBody: `{"error":"error to calculate path: there are no flights"}`,
		},
		{
			name: "invalid itinerary data",
			fields: fields{
				responseWriter: httptest.NewRecorder(),
			},
			args: args{
				err:     domain.ErrInvalidItinerary,
				details: "error to calculate path",
			},
			wantStatusCode:   422,
			wantResponseBody: `{"error":"error to calculate path: invalid itinerary data"}`,
		},
		{
			name: "unknown error",
			fields: fields{
				responseWriter: httptest.NewRecorder(),
			},
			args: args{
				err:     errors.New("any error"),
				details: "error to calculate path",
			},
			wantStatusCode:   503,
			wantResponseBody: `{"error":"error to calculate path: any error"}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			output := jsonOutput{w: tt.fields.responseWriter}

			if err := output.domainError(tt.args.err, tt.args.details); err != nil {
				t.Fatalf(err.Error())
			}

			responseWriter := tt.fields.responseWriter

			assertHTTPResponse(t, responseWriter.Result(), tt.wantStatusCode, tt.wantResponseBody)
		})
	}
}

func assertHTTPResponse(t *testing.T, response *http.Response, expectedStatusCode int, expectedPayload string) {
	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("error to read payload response: %v", err)
	}

	if string(rawResponse) != expectedPayload {
		t.Errorf("Payload response does not match, got=%s expected=%s", string(rawResponse), expectedPayload)
	}

	if response.StatusCode != expectedStatusCode {
		t.Errorf("Status code does not match, got=%d expected=%d", response.StatusCode, expectedStatusCode)
	}
}
