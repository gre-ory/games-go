package api

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
)

func (s *gameServer) htmx_leave_game(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("[api] htmx_leave_game", zap.String("path", r.URL.Path))

	ctx := r.Context()

	var cookie *Cookie
	var game *model.Game
	var err error

	switch {
	default:

		cookie, err = s.GetCookie(r)
		if err != nil {
			s.logger.Info("cookie not found!", zap.Error(err))
			break
		}
		err = cookie.Validate()
		if err != nil {
			break
		}

		gameId, err := s.extractPathGameId(ctx)
		if err != nil {
			break
		}

		game, err = s.service.LeaveGame(gameId, cookie.PlayerId)
		if err != nil {
			break
		}

		// s.hub.UnsetGameId(cookie.PlayerId)

		s.broadcastGame(game)
		s.broadcastJoinableGames()

		return
	}

	// error response

	s.renderError(w, err)
}
