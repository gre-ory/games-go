package model

import (
	"github.com/jaevor/go-nanoid"

	"github.com/gre-ory/games-go/internal/util"
)

type GameId string

var gameIdAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var generateGameId = util.Must(nanoid.CustomASCII(gameIdAlphabet, 4))

func NewGameId() GameId {
	return GameId(generateGameId())
}
