package model

import "fmt"

var (
	ErrMissingUserId       = fmt.Errorf("missing user id")
	ErrMissingUserName     = fmt.Errorf("missing user name")
	ErrMissingUserAvatar   = fmt.Errorf("missing user avatar")
	ErrInvalidUserAvatar   = fmt.Errorf("invalid user avatar")
	ErrMissingUserLanguage = fmt.Errorf("missing user language")
	ErrInvalidUserLanguage = fmt.Errorf("invalid user language")
)
