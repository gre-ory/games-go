package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/game/user/model"
)

// //////////////////////////////////////////////////
// create session

func (s *sessionServer) session_create(w http.ResponseWriter, r *http.Request) {

	var session *model.Session
	var err error

	switch {
	default:

		session, err = s.service.Create()
		if err != nil {
			break
		}

		// success response

		s.SetSessionHeaders(w, session)
		return
	}

	// error response

	// TODO better response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
