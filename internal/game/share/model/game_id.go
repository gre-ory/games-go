package model

import (
	"github.com/jaevor/go-nanoid"

	"github.com/gre-ory/games-go/internal/util"
)

type GameId string

const gameIdAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var gameIdGenerateFn = util.Must(nanoid.CustomASCII(gameIdAlphabet, 6))

var GenerateGameId = func() GameId {
	return GameId(gameIdGenerateFn())
}
