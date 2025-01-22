package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Server struct {
	Echo   *echo.Echo
	Server *http.Server
}

func NewServer(e *echo.Echo, addr string) *Server {
	return &Server{
		Echo: e,
		Server: &http.Server{
			Addr:    addr,
			Handler: e,
		},
	}
}

func (s *Server) Start() {
	go func() {
		log.Println("Starting server on", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) {
	log.Println("Shutting down server...")
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}
	log.Println("Server gracefully stopped")
}

func (s *Server) WaitForShutdownSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
}
