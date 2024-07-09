package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/util"

	"github.com/gre-ory/games-go/internal/game/czm/model"
)

func (s *gameServer) onMessage(playerId share_model.PlayerId, message []byte) {

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
			err = share_model.ErrMissingAction
			break
		}

		player, err = s.GetPlayer(playerId)
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
			gameId := share_model.GameId(jsonMessage.GameId)
			err = s.HandleJoinGame(player, gameId)
		case "start-game":
			err = s.HandleStartGame(player)
		case "select-card":
			err = s.HandleSelectCard(player, jsonMessage)
		case "play-card":
			err = s.HandlePlayCard(player, jsonMessage)
		case "leave-game":
			err = s.HandleLeaveGame(player)
		default:
			err = model.ErrMissingAction
		}
	}

	if err != nil {
		s.BroadcastErrorToPlayer(playerId, err)
	}
}

type JsonMessage struct {
	// Headers    *JsonHeaders `json:"HEADERS,omitempty"`
	Action           string `json:"action,omitempty"`
	PlayerName       string `json:"name,omitempty"`
	GameId           string `json:"game,omitempty"`
	CardNumberStr    string `json:"card,omitempty"`
	DiscardNumberStr string `json:"discard,omitempty"`
}

func (j *JsonMessage) CardNumber() (int, error) {
	if j.CardNumberStr == "" {
		return 0, model.ErrInvalidCardNumber
	}
	cardNumber := util.ToInt(j.CardNumberStr)
	if cardNumber == 0 {
		return 0, model.ErrInvalidCardNumber
	}
	return cardNumber, nil
}

func (j *JsonMessage) DiscardNumber() (int, error) {
	if j.DiscardNumberStr == "" {
		return 0, model.ErrInvalidDiscardNumber
	}
	discardNumber := util.ToInt(j.DiscardNumberStr)
	if discardNumber == 0 {
		return 0, model.ErrInvalidDiscardNumber
	}
	return discardNumber, nil
}

type JsonHeaders struct {
	HxRequest     string `json:"HX-Request,omitempty"`
	HxTrigger     string `json:"HX-Trigger,omitempty"`
	HxTriggerName string `json:"HX-Trigger-Name,omitempty"`
	HxTarget      string `json:"HX-Target,omitempty"`
	HxCurrentUrl  string `json:"HX-Current-URL,omitempty"`
}
