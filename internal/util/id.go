package util

import (
	"github.com/jaevor/go-nanoid"
)

var gameIdAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var GenerateGameId = Must(nanoid.CustomASCII(gameIdAlphabet, 4))

var playerIdAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var GeneratePlayerId = Must(nanoid.CustomASCII(playerIdAlphabet, 6))

var tokenAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
var GenerateUserToken = Must(nanoid.CustomASCII(tokenAlphabet, 10))
