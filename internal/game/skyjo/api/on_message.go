package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/skyjo/model"
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
		case "create-game":
			err = s.HandleCreateGame(player)
		case "join-game":
			err = s.HandleJoinGame(player)
		case "start-game":
			err = s.HandleStartGame(player)
		case "play":
			err = s.HandlePlay(player, jsonMessage)
		case "leave-game":
			err = s.HandleLeaveGame(player)
		default:
			err = model.ErrMissingAction
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
