package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"

	share_model "github.com/gre-ory/games-go/internal/game/share/model"

	"github.com/gre-ory/games-go/internal/game/skj/model"
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
		// case "set-name":
		// 	err = s.ws_set_player_name(player, jsonMessage)
		case "create-game":
			err = s.HandleCreateGame(player)
		case "join-game":
			gameId := share_model.GameId(jsonMessage.GameId)
			err = s.HandleJoinGame(player, gameId)
		case "start-game":
			err = s.HandleStartGame(player)
		case "draw-discard-card":
			err = s.HandleDrawDiscardCard(player)
		case "draw-card":
			err = s.HandleDrawCard(player)
		case "put-card":
			err = s.HandlePutCard(player, jsonMessage)
		case "discard-card":
			err = s.HandleDiscardCard(player)
		case "flip-card":
			err = s.HandleFlipCard(player, jsonMessage)
		case "leave-game":
			err = s.HandleLeaveGame(player)
		default:
			err = share_model.ErrMissingAction
		}
	}

	if err != nil {
		s.BroadcastErrorToPlayer(playerId, err)
	}
}

type JsonMessage struct {
	// Headers    *JsonHeaders `json:"HEADERS,omitempty"`
	Action          string `json:"action,omitempty"`
	PlayerName      string `json:"name,omitempty"`
	GameId          string `json:"game,omitempty"`
	ColumnNumberStr string `json:"column,omitempty"`
	RowNumberStr    string `json:"row,omitempty"`
}

func (j *JsonMessage) ColumnNumber() (int, error) {
	if j.ColumnNumberStr == "" {
		return 0, model.ErrInvalidColumn
	}
	columnNumber := util.ToInt(j.ColumnNumberStr)
	if columnNumber == 0 {
		return 0, model.ErrInvalidColumn
	}
	return columnNumber, nil
}

func (j *JsonMessage) RowNumber() (int, error) {
	if j.RowNumberStr == "" {
		return 0, model.ErrInvalidRow
	}
	rowNumber := util.ToInt(j.RowNumberStr)
	if rowNumber == 0 {
		return 0, model.ErrInvalidRow
	}
	return rowNumber, nil
}

func (j *JsonMessage) Cell() (int, int, error) {
	columnNumber, err := j.ColumnNumber()
	if err != nil {
		return 0, 0, err
	}
	rowNumber, err := j.RowNumber()
	if err != nil {
		return 0, 0, err
	}
	return columnNumber, rowNumber, nil
}

type JsonHeaders struct {
	HxRequest     string `json:"HX-Request,omitempty"`
	HxTrigger     string `json:"HX-Trigger,omitempty"`
	HxTriggerName string `json:"HX-Trigger-Name,omitempty"`
	HxTarget      string `json:"HX-Target,omitempty"`
	HxCurrentUrl  string `json:"HX-Current-URL,omitempty"`
}
