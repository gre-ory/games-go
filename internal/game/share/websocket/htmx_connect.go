package websocket

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/util"
)

// //////////////////////////////////////////////////
// htmx connect

func (s *hubServer[PlayerT, GameT]) HtmxConnect(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("[api] htmx_connect ", zap.String("path", r.URL.Path))

	var cookie *model.Cookie
	var user User
	var game GameT
	var found bool
	var err error

	switch {
	default:

		//
		// extract cookie
		//

		cookie, err = s.cookierServer.GetValidCookie(r)
		if err != nil {
			s.logger.Info("[api] no valid cookie >>> STOP", zap.Error(err))
			break
		}
		userId := cookie.Id

		//
		// fetch ( or create ) websocket user
		//

		s.logger.Info(fmt.Sprintf("[api] cookie %s >>> getting user...", userId), zap.Any("cookie", cookie))
		user, err = s.Hub().GetUser(userId)
		if err != nil {
			if !errors.Is(err, model.ErrUserNotFound) {
				s.logger.Info(fmt.Sprintf("[api] user %s not found >>> ERROR", userId), zap.Error(err))
				break
			}
			s.logger.Info(fmt.Sprintf("[api] user %s not found >>> create a new one", userId))
			user = s.newUserFromCookieFn(cookie)
			s.RegisterUser(user)
		} else {
			s.logger.Info(fmt.Sprintf("[api] user %s already exists", userId), zap.Any("user", user))
		}

		//
		// connect socket
		//

		s.logger.Info(fmt.Sprintf("[api] user %s >>> connecting...", userId))
		err = user.ConnectSocket(w, r)
		if err != nil {
			s.logger.Info(fmt.Sprintf("[api] user %s >>> connection failed", userId), zap.Error(err))
			break
		}
		s.logger.Info(fmt.Sprintf("[api] ... user %s connected", userId))

		//
		// broadcast joinable games ( if not playing )
		//

		if !user.HasGameId() {
			s.logger.Info(fmt.Sprintf("[api] user %s >>> broadcasting games...", userId))
			s.BroadcastJoinableGamesToUser(userId)
			return
		}

		//
		// broadcast game layout to player ( if playing )
		//

		gameId := user.GameId()
		s.logger.Info(fmt.Sprintf("[api] user %s >>> fetching game %s...", userId, gameId))
		game, err = s.service.GetGame(gameId)
		if err != nil {
			break
		}

		playerId := user.PlayerId()
		_, found = game.Player(playerId)
		if !found {
			s.logger.Info(fmt.Sprintf("[api] user %s >>> not found in game %s", userId, gameId))
			err = model.ErrPlayerNotFound
			break
		}

		s.logger.Info(fmt.Sprintf("[api] player %s >>> broadcasting game layout...", playerId))
		s.BroadcastGameLayoutToPlayer(playerId, game)

		//
		// broadcast game to other players
		//

		s.logger.Info(fmt.Sprintf("[api] player %s >>> broadcasting game...", playerId))
		s.BroadcastGame(game)

		return

	}

	// error response
	s.logger.Info("[api] htmx_connect: FAILED", zap.String("path", r.URL.Path), zap.Error(err))
	util.EncodeJsonErrorResponse(w, err)
}
