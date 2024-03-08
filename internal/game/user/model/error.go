package model

import "fmt"

var (
	ErrSessionNotFound      = fmt.Errorf("session not found")
	ErrSessionExpired       = fmt.Errorf("session expired")
	ErrSessionTokenNotFound = fmt.Errorf("session token not found")
)
