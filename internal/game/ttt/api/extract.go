package api

import (
	"context"

	"github.com/gre-ory/games-go/internal/game/ttt/model"
	"github.com/gre-ory/games-go/internal/util"
)

func (s *gameServer) extractPathGameId(ctx context.Context) (model.GameId, error) {
	gameId := model.GameId(util.ExtractPathParameter(ctx, "game_id"))
	if gameId == "" {
		return "", model.ErrMissingGameId
	}
	return gameId, nil
}

func (s *gameServer) extractPathPlayX(ctx context.Context) (int, error) {
	playX := util.ExtractPathIntParameter(ctx, "play_x")
	if playX == 0 {
		return 0, model.ErrMissingPlayX
	}
	return playX, nil
}

func (s *gameServer) extractPathPlayY(ctx context.Context) (int, error) {
	playY := util.ExtractPathIntParameter(ctx, "play_y")
	if playY == 0 {
		return 0, model.ErrMissingPlayY
	}
	return playY, nil
}
