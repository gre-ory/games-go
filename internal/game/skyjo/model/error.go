package model

import "fmt"

var (
	ErrGameNotFound         = fmt.Errorf("game not found")
	ErrPlayerNotFound       = fmt.Errorf("player not found")
	ErrPlayerAlreadyPlaying = fmt.Errorf("player already playing")
	ErrMissingGameId        = fmt.Errorf("missing game id")
	ErrMissingPlayerId      = fmt.Errorf("missing player id")
	ErrMissingPlayerName    = fmt.Errorf("missing player name")
	ErrMissingPlayers       = fmt.Errorf("missing players")
	ErrGameNotStarted       = fmt.Errorf("game not started")
	ErrGameAlreadyStarted   = fmt.Errorf("game already started")
	ErrGameNotJoinable      = fmt.Errorf("game not joinable")
	ErrOutOfColumnBound     = fmt.Errorf("out of column bound")
	ErrOutOfRowBound        = fmt.Errorf("out of row bound")
	ErrWrongPlayer          = fmt.Errorf("wrong player")
	ErrMissingPlayX         = fmt.Errorf("missing play x")
	ErrMissingPlayY         = fmt.Errorf("missing play y")
	ErrGameStopped          = fmt.Errorf("game stopped")
	ErrGameNotStopped       = fmt.Errorf("game not stopped")
	ErrAlreadyPlayOnCell    = fmt.Errorf("already played on cell")
	ErrMissingAction        = fmt.Errorf("missing action")
	ErrUnknownAction        = fmt.Errorf("unknown action")
	ErrPlayerAlreadyActive  = fmt.Errorf("player already active")
)
