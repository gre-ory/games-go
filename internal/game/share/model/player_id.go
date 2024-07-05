package model

import (
	"github.com/jaevor/go-nanoid"

	"github.com/gre-ory/games-go/internal/util"
)

type PlayerId string

const playerIdAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var playerIdGenerateFn = util.Must(nanoid.CustomASCII(playerIdAlphabet, 6))

var GeneratePlayerId = func() PlayerId {
	return PlayerId(playerIdGenerateFn())
}
