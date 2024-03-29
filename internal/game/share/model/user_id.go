package model

import (
	"github.com/jaevor/go-nanoid"

	"github.com/gre-ory/games-go/internal/util"
)

type UserId string

var userIdAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var generateUserId = util.Must(nanoid.CustomASCII(userIdAlphabet, 6))

func NewUserId() UserId {
	return UserId(generateUserId())
}

func (i UserId) Validate() error {
	if i == "" {
		return ErrMissingUserId
	}
	return nil
}
