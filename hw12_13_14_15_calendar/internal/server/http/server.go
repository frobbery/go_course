package internalhttp

import (
	"context"
	"github.com/gorilla/mux"
)

type Server struct {
	host string
	port string
}

type Logger interface {
	Info(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, host string, port string) *Server {
	return &Server{host: host, port: port}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}
