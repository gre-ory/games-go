package service

import (
	"github.com/gre-ory/games-go/internal/game/user/model"
	"github.com/gre-ory/games-go/internal/game/user/store"
)

type SessionService interface {
	Create() (*model.Session, error)
	Retrieve(token model.Token) (*model.Session, error)
	AttachContext(token model.Token, context map[string]string) (*model.Session, error)
}

func NewSessionService(store store.SessionStore) SessionService {
	return &sessionService{
		store: store,
	}
}

type sessionService struct {
	store store.SessionStore
}

func (s *sessionService) Create() (*model.Session, error) {
	session := model.NewSession()
	err := s.store.Set(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *sessionService) Retrieve(token model.Token) (*model.Session, error) {
	session, err := s.store.Get(token)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		s.store.Delete(token)
		return nil, model.ErrSessionExpired
	}
	return session, nil
}

func (s *sessionService) AttachContext(token model.Token, context map[string]string) (*model.Session, error) {
	session, err := s.store.Get(token)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		s.store.Delete(token)
		return nil, model.ErrSessionExpired
	}
	for key, value := range context {
		session.Context[key] = value
	}
	err = s.store.Set(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}
