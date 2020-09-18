package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/MattIzSpooky/tf2.rest/responses"
)

const randomResponseRoute = "/random"

func (s *Server) randomResponseHandler(w http.ResponseWriter, _ *http.Request) {
	randomResponse := responses.All[rand.Intn(len(responses.All))]
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
