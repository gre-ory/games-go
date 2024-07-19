package model

import (
	"fmt"
	"strings"
)

type PlayerId string

func NewPlayerId(gameId GameId, userId UserId) PlayerId {
	if gameId == "" {
		panic(ErrMissingGameId)
	}
	if userId == "" {
		panic(ErrMissingUserId)
	}
	return PlayerId(fmt.Sprintf("%s-%s", gameId, userId))
}

func (id PlayerId) GameId() GameId {
	parts := strings.Split(string(id), "-")
	if len(parts) != 2 {
		panic(ErrInvalidPlayerId)
	}
	return GameId(parts[0])
}

func (id PlayerId) UserId() UserId {
	parts := strings.Split(string(id), "-")
	if len(parts) != 2 {
		panic(ErrInvalidPlayerId)
	}
	return UserId(parts[1])
}

func (id PlayerId) MatchUser(userId UserId) bool {
	return id.UserId() == userId
}
