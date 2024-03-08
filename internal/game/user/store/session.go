package store

import "github.com/gre-ory/games-go/internal/game/user/model"

type SessionStore interface {
	Set(session *model.Session) error
	Get(token model.Token) (*model.Session, error)
	Delete(token model.Token) error
}

func NewSessionStore() SessionStore {
	return &sessionStore{
		sessions: map[model.Token]*model.Session{},
	}
}

type sessionStore struct {
	sessions map[model.Token]*model.Session
}

func (s *sessionStore) Set(session *model.Session) error {
	s.sessions[session.Token] = session
	return nil
}

func (s *sessionStore) Get(token model.Token) (*model.Session, error) {
	if session, ok := s.sessions[token]; ok {
		return session, nil
	}
	return nil, model.ErrSessionNotFound
}

func (s *sessionStore) Delete(token model.Token) error {
	if _, ok := s.sessions[token]; ok {
		delete(s.sessions, token)
		return nil
	}
	return model.ErrSessionNotFound
}
