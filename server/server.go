package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/rs/cors"
)

// Server is our server object. A normal Server object with a logger attached to it.
type Server struct {
	http.Server
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
}

// NewServer creates an instance of Server
func NewServer() *Server {
	server := &Server{
		Server: http.Server{
			Addr:         fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		infoLogger:  log.New(os.Stdout, color.Blue.Sprint("[Info]: "), log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, color.Red.Sprint("[Error]: "), log.Ldate|log.Ltime),
		warnLogger:  log.New(os.Stdout, color.Yellow.Sprint("[Warn]: "), log.Ldate|log.Ltime),
	}

	server.addHandlers()

	return server
}

func (s *Server) addHandlers() {
	mux := &http.ServeMux{}

	mux.Handle("/", jsonMiddleware(http.HandlerFunc(s.randomQuoteHandler)))

	s.Handler = cors.Default().Handler(mux)
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

// Serve makes our endpoints available
func (s *Server) Serve() error {
	s.infoLogger.Println(fmt.Sprintf("Running HTTP server on: http://%s", s.Addr))

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.errorLogger.Printf("Could not listen on %s: %v\n", s.Addr, err.Error())

		return err
	}

	return nil
}

// GracefulShutdown closes our server gracefully, waiting for requests to be handled before closing down.
func (s *Server) GracefulShutdown() error {
	s.infoLogger.Println("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)

	if err := s.Shutdown(ctx); err != nil {
		s.errorLogger.Printf("Could not gracefully shutdown the server: %v\n", err.Error())
		return err
	}

	s.infoLogger.Println("Server has successfully closed down")

	return nil
}
