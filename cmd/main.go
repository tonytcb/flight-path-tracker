package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tonytcb/flight-path-tracker/pkg/api/http"
	"github.com/tonytcb/flight-path-tracker/pkg/infra/flightparser"
	"github.com/tonytcb/flight-path-tracker/pkg/usecase"
)

func main() {
	const httpPort = 8080

	log.Println("Starting application")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	var (
		flightsCalculatorHandler = http.NewFlightCalculatorHandler(
			flightparser.NewJSONParser(),
			usecase.NewFlightTracker(),
		)
		httpServer = http.NewServer(
			flightsCalculatorHandler,
		)
	)

	if err := httpServer.Start(httpPort); err != nil {
		log.Fatalf(err.Error())
	}

	<-done

	if err := httpServer.Stop(context.Background()); err != nil {
		log.Fatalf("error to shutdown http server: %v", err)
	}

	log.Println("Shutting down application")
}
