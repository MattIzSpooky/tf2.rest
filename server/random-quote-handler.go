package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/MattIzSpooky/tf2.rest/quotes"
)

func (s *Server) randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	randomQuote := quotes.ALL[rand.Intn(len(quotes.ALL))]
	json, err := json.Marshal(randomQuote)

	if err != nil {
		s.errorLogger.Println(fmt.Sprintf("Could not marshal quote: %s", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(json)

	if err != nil {
		s.errorLogger.Println(fmt.Sprintf("Could not write JSON: %s", err.Error()))
	}
}
