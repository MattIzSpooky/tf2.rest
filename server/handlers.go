package server

import (
	"github.com/rs/cors"
	"net/http"
)

const rootRoute = "/"

func (s *Server) addHandlers() {
	mux := &http.ServeMux{}

	mux.Handle(rootRoute, jsonMiddleware(http.HandlerFunc(s.randomResponseHandler)))
	mux.Handle(byClassRoute, jsonMiddleware(http.HandlerFunc(s.randomByClassHandler)))

	s.Handler = cors.Default().Handler(mux)
}
