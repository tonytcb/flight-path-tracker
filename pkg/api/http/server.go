package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type Server struct {
	httpServer *http.Server

	calculatorHandler *FlightCalculatorHandler
}

func NewServer(
	calculatorHandler *FlightCalculatorHandler,
) *Server {
	return &Server{
		calculatorHandler: calculatorHandler,
	}
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.healthHandler)
	mux.HandleFunc("/calculate", s.calculatorHandler.Handle)

	log.Println("Starting HTTP Server on port", port)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error to listen and serve http api: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Shutting down HTTP Server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "error to shutdown http server")
	}

	return nil
}

func (s *Server) healthHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "OK")
}
