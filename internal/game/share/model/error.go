package model

import "fmt"

var (
	ErrInvalidCookie         = fmt.Errorf("invalid cookie")
	ErrMissingUserId         = fmt.Errorf("missing user id")
	ErrMissingUserName       = fmt.Errorf("missing user name")
	ErrMissingUserAvatar     = fmt.Errorf("missing user avatar")
	ErrInvalidUserAvatar     = fmt.Errorf("invalid user avatar")
	ErrUnsupportedLanguage   = fmt.Errorf("unsupported language")
	ErrMissingGameId         = fmt.Errorf("missing game id")
	ErrMissingPlayers        = fmt.Errorf("missing players")
	ErrGameNotFound          = fmt.Errorf("game not found")
	ErrPlayerNotFound        = fmt.Errorf("player not found")
	ErrPlayerNotInGame       = fmt.Errorf("player not in game")
	ErrWrongPlayer           = fmt.Errorf("wrong player")
	ErrGameNotJoinable       = fmt.Errorf("game not joinable")
	ErrGameNotStarted        = fmt.Errorf("game not started")
	ErrGameAlreadyStarted    = fmt.Errorf("game already started")
	ErrGameStopped           = fmt.Errorf("game stopped")
	ErrGameNotStopped        = fmt.Errorf("game not stopped")
	ErrGameNotStartable      = fmt.Errorf("game not startable")
	ErrGameMarkedForDeletion = fmt.Errorf("game marked for deletion")
)
