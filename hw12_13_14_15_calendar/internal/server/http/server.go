package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	//nolint:depguard
	"github.com/gorilla/mux"
)

type Server struct {
	host    string
	port    string
	logger  Logger
	server  *http.Server
	logFile *os.File
}

type Logger interface {
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, _ Application, host string, port string) *Server {
	return &Server{host: host, port: port, logger: logger}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/", s.helloWorldRequest)

	err := s.initLogFile()
	if err != nil {
		return err
	}
	handler := loggingMiddleware(router)

	//nolint:gosec
	s.server = &http.Server{
		Addr:    s.host + ":" + s.port,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()

	return s.server.ListenAndServe()
}

func (s *Server) initLogFile() error {
	logFile, err := os.OpenFile("calendar.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		s.logger.Error("could not init log: " + err.Error())
		return err
	}
	s.logFile = logFile
	log.SetOutput(logFile)
	return nil
}

func (s *Server) helloWorldRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, World!")
}

func (s *Server) Stop(ctx context.Context) error {
	ctxShutdown, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctxShutdown)
	if err != nil {
		s.logger.Error("could not shut down server: " + err.Error())
	}
	s.logFile.Close()
	return err
}
