package model

import (
	"time"

	"github.com/gre-ory/games-go/internal/util"
)

type Session struct {
	Token   Token
	Expire  time.Time
	Context map[string]string
}

func NewSession() *Session {
	return &Session{
		Token:   Token(util.GenerateUserToken()),
		Expire:  time.Now().Add(24 * time.Hour),
		Context: map[string]string{},
	}
}

func (t *Session) IsExpired() bool {
	return t.Expire.Before(time.Now())
}
