package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	share_websocket "github.com/gre-ory/games-go/internal/game/share/websocket"

	"github.com/gre-ory/games-go/internal/game/czm/model"
)

func (s *gameServer) onMessage(userId share_model.UserId, message []byte) {

	var jsonMessage JsonMessage
	var user share_websocket.User
	var player *model.Player
	var err error

	switch {
	default:

		//
		// decode message
		//

		err = json.NewDecoder(bytes.NewReader(message)).Decode(&jsonMessage)
		if err != nil {
			s.logger.Info("[DEBUG] failed to decode message: " + err.Error())
			break
		}

		if jsonMessage.Action == "" {
			err = share_model.ErrMissingAction
			break
		}

		//
		// fetch websocket user
		//

		if userId == "" {
			err = share_model.ErrMissingUserId
			break
		}
		user, err = s.GetUser(userId)
		if err != nil {
			break
		}
		if user.IsInactive() {
			err = share_model.ErrInactiveUser
			break
		}

		//
		// create or join game ( if not playing )
		//

		if !user.HasGameId() {

			now := time.Now()
			s.logger.Info(fmt.Sprintf("[WS] ------------------------- user %s :: %s -------------------------", userId, jsonMessage.Action), zap.Any("message", jsonMessage))
			defer func() {
				s.logger.Info(fmt.Sprintf("[WS] ------------------------- user %s :: %s ( %s ) -------------------------", userId, jsonMessage.Action, time.Since(now)))
			}()

			switch jsonMessage.Action {
			case "create-game":
				err = s.HandleCreateGame(user)
			case "join-game":
				err = s.HandleJoinGame(jsonMessage.GameId(), user)
			default:
				err = share_model.ErrInvalidAction
			}
			break
		}

		//
		// fetch player ( if playing )
		//

		playerId := user.PlayerId()
		player, err = s.GetPlayer(playerId)
		if err != nil {
			break
		}

		now := time.Now()
		s.logger.Info(fmt.Sprintf("[WS] ------------------------- player %s :: %s -------------------------", playerId, jsonMessage.Action), zap.Any("message", jsonMessage))
		defer func() {
			s.logger.Info(fmt.Sprintf("[WS] ------------------------- player %s :: %s ( %s ) -------------------------", playerId, jsonMessage.Action, time.Since(now)))
		}()

		switch jsonMessage.Action {
		case "start-game":
			err = s.HandleStartGame(player)
		case "select-card":
			err = s.HandleSelectCard(player, jsonMessage.CardNumber())
		case "play-card":
			err = s.HandlePlayCard(player, jsonMessage.DiscardNumber())
		case "leave-game":
			err = s.HandleLeaveGame(player)
		default:
			err = share_model.ErrInvalidAction
		}
	}

	if userId != "" && err != nil {
		s.BroadcastErrorToUser(userId, err)
	}
}

type JsonMessage struct {
	// Headers    *JsonHeaders `json:"HEADERS,omitempty"`
	Action           string `json:"action,omitempty"`
	PlayerName       string `json:"name,omitempty"`
	GameIdStr        string `json:"game,omitempty"`
	CardNumberStr    string `json:"card,omitempty"`
	DiscardNumberStr string `json:"discard,omitempty"`
}

func (j *JsonMessage) GameId() share_model.GameId {
	return share_model.GameId(j.GameIdStr)
}

func (j *JsonMessage) CardNumber() int {
	return util.ToInt(j.CardNumberStr)
}

func (j *JsonMessage) DiscardNumber() int {
	return util.ToInt(j.DiscardNumberStr)
}

type JsonHeaders struct {
	HxRequest     string `json:"HX-Request,omitempty"`
	HxTrigger     string `json:"HX-Trigger,omitempty"`
	HxTriggerName string `json:"HX-Trigger-Name,omitempty"`
	HxTarget      string `json:"HX-Target,omitempty"`
	HxCurrentUrl  string `json:"HX-Current-URL,omitempty"`
}
