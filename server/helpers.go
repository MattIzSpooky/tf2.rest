package server

import (
	"fmt"
	"net/http"
)

func (s *Server) writeError(w http.ResponseWriter, text string, statusCode int)  {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(fmt.Sprintf(`{ "error": "%s" }`, text)))

	if err != nil {
		s.errorLogger.Println(fmt.Sprintf("Could not write error: %s", err.Error()))
	}
}

func (s *Server) crash(w http.ResponseWriter, reason string) {
	s.errorLogger.Println(reason)
	w.WriteHeader(http.StatusInternalServerError)
}
