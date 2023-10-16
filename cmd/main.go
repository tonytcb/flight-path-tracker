package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/tonytcb/flight-path-tracker/pkg/api/http"
	"github.com/tonytcb/flight-path-tracker/pkg/infra/flightparser"
	"github.com/tonytcb/flight-path-tracker/pkg/usecase"
)

const (
	httpPortEnVarName = "HTTP_PORT"
	httpPortDefault   = 8080
)

func main() {
	log.Println("Starting application")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	httpPort, err := loadEnvVarInt(httpPortEnVarName, httpPortDefault)
	if err != nil {
		log.Fatalf("error to load env var %s: %v", httpPortEnVarName, err)
	}

	/**
	 * To have exactly the same input api provided in the examples (json containing a list of arrays),
	 * it's easily done change injecting the flightparser.NewJSONOfArraysParser() instead of flightparser.NewJSONParser().
	 */

	var (
		flightsCalculatorHandler = http.NewFlightCalculatorHandler(
			flightparser.NewJSONParser(),
			usecase.NewFlightTracker(),
		)
		httpServer = http.NewServer(
			flightsCalculatorHandler,
		)
	)

	if err = httpServer.Start(httpPort); err != nil {
		log.Fatalf(err.Error())
	}

	<-done

	if err := httpServer.Stop(context.Background()); err != nil {
		log.Fatalf("error to shutdown http server: %v", err)
	}

	log.Println("Shutting down application")
}

func loadEnvVarInt(keyName string, defaultValue int) (int, error) {
	if v := os.Getenv(keyName); v != "" {
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}

		return intValue, nil
	}

	return defaultValue, nil
}
