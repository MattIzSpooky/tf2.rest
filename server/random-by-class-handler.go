package server

import (
	"encoding/json"
	"fmt"
	"github.com/MattIzSpooky/tf2.rest/class"
	"github.com/MattIzSpooky/tf2.rest/responses"
	"math/rand"
	"net/http"
	"strings"
)

const byClassRoute = rootRoute + "by-class/"

func (s *Server) randomByClassHandler(w http.ResponseWriter, r *http.Request) {
	classFromReq := strings.TrimPrefix(r.URL.Path, byClassRoute)

	if classFromReq == "" {
		s.writeError(w, fmt.Sprintf("No class was provided in route. Example: %sscout", byClassRoute), http.StatusBadRequest)
		return
	}

	if !class.Contains(classFromReq) {
		s.writeError(w, fmt.Sprintf("Class does not exist."), http.StatusBadRequest)
		return
	}

	classResponse := responses.FilterByClass(classFromReq)

	randomResponse := classResponse[rand.Intn(len(classResponse))]
	response, err := json.Marshal(randomResponse)

	if err != nil {
		s.crash(w, fmt.Sprintf("Could not marshal response: %s", err.Error()))
		return
	}

	_, err = w.Write(response)

	if err != nil {
		s.crash(w, fmt.Sprintf("Could not write JSON: %s", err.Error()))
	}
}
