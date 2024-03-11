package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
	"go.uber.org/zap"
)

func (s *gameServer) onMessage(playerId model.PlayerId, message []byte) {

	var jsonMessage JsonMessage
	var player *model.Player
	var err error

	switch {
	default:

		err = json.NewDecoder(bytes.NewReader(message)).Decode(&jsonMessage)
		if err != nil {
			s.logger.Info("[DEBUG] failed to decode message: " + err.Error())
			break
		}

		if jsonMessage.Action == "" {
			err = model.ErrMissingAction
			break
		}

		player, err = s.hub.GetPlayer(playerId)
		if err != nil {
			s.logger.Info("[DEBUG] player not founf: " + err.Error())
			break
		}

		now := time.Now()
		s.logger.Info(fmt.Sprintf("[WS] ------------------------- %s :: %s -------------------------", playerId, jsonMessage.Action), zap.Any("message", jsonMessage))
		defer func() {
			s.logger.Info(fmt.Sprintf("[WS] ------------------------- %s :: %s ( %s ) -------------------------", playerId, jsonMessage.Action, time.Since(now)))
		}()

		switch jsonMessage.Action {
		case "set-name":
			err = s.ws_set_player_name(player, jsonMessage)
		case "create-game":
			err = s.ws_create_game(player, jsonMessage)
		case "join-game":
			err = s.ws_join_game(player, jsonMessage)
		case "play":
			err = s.ws_play(player, jsonMessage)
		}
	}

	if err != nil {
		s.broadcastErrorToPlayer(playerId, err)
	}
}

type JsonMessage struct {
	// Headers    *JsonHeaders `json:"HEADERS,omitempty"`
	Action     string `json:"action,omitempty"`
	PlayerName string `json:"name,omitempty"`
	GameId     string `json:"game,omitempty"`
	PlayX      string `json:"x,omitempty"`
	PlayY      string `json:"y,omitempty"`
}

type JsonHeaders struct {
	HxRequest     string `json:"HX-Request,omitempty"`
	HxTrigger     string `json:"HX-Trigger,omitempty"`
	HxTriggerName string `json:"HX-Trigger-Name,omitempty"`
	HxTarget      string `json:"HX-Target,omitempty"`
	HxCurrentUrl  string `json:"HX-Current-URL,omitempty"`
}
