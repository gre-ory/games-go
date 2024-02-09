package service

import (
	"github.com/gre-ory/games-go/internal/game/tictactoe/model"
)

var (
	games = make(map[model.GameId]*model.Game)
)

func NewGame(playerName string) (*model.Game, *model.Player) {
	player := model.NewPlayer(playerName)
	game := model.NewGame().WithPlayer(player)
	games[game.Id] = game
	return game, player
}

func JoinGame(gameId model.GameId, playerName string) (*model.Game, *model.Player) {
	game := games[gameId]
	player := model.NewPlayer(playerName)
	return game.WithPlayer(player), player
}

func DeleteGame(gameId model.GameId) *model.Game {
	game := games[gameId]
	delete(games, gameId)
	return game
}
